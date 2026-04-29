package notification

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// RenderTemplate 极简 {{key}} 占位替换。
func RenderTemplate(tpl string, vars map[string]string) string {
	out := tpl
	for k, v := range vars {
		out = strings.ReplaceAll(out, "{{"+k+"}}", v)
	}
	return out
}

// EmailSender SMTP 发送实现（读取当前启用邮件渠道配置）。
type EmailSender struct{}

func (EmailSender) Send(ctx context.Context, to, subject, body string) error {
	provider, cfg, ok := GetActiveEmailConfig(ctx)
	if !ok {
		return gerror.New("邮件渠道未配置或当前启用配置无效")
	}
	switch provider {
	case "smtp":
		smtpCfg, ok := cfg.(*SMTPConfig)
		if !ok || smtpCfg == nil {
			return gerror.New("SMTP 配置解析失败")
		}
		return sendBySMTP(ctx, smtpCfg, to, subject, body)
	case "sendgrid":
		sendgridCfg, ok := cfg.(*SendGridConfig)
		if !ok || sendgridCfg == nil {
			return gerror.New("SendGrid 配置解析失败")
		}
		return sendBySendGrid(ctx, sendgridCfg, to, subject, body)
	default:
		return gerror.Newf("暂不支持的邮件提供商: %s", provider)
	}
}

func (EmailSender) SendTemplate(ctx context.Context, to, subject, templateID string, params map[string]string) error {
	provider, cfg, ok := GetActiveEmailConfig(ctx)
	if !ok {
		return gerror.New("邮件渠道未配置或当前启用配置无效")
	}
	if provider != "sendgrid" {
		return gerror.New("模板邮件仅支持 SendGrid 提供商")
	}
	sendgridCfg, ok := cfg.(*SendGridConfig)
	if !ok || sendgridCfg == nil {
		return gerror.New("SendGrid 配置解析失败")
	}
	return sendTemplateBySendGrid(ctx, sendgridCfg, to, subject, templateID, params)
}

func sendBySMTP(ctx context.Context, smtpCfg *SMTPConfig, to, subject, body string) error {
	if strings.TrimSpace(smtpCfg.Host) == "" || smtpCfg.Port <= 0 {
		return gerror.New("SMTP 配置缺少 host/port")
	}
	if strings.TrimSpace(smtpCfg.From) == "" {
		return gerror.New("SMTP 配置缺少发件人 from")
	}

	fromAddr, err := parseAddressMailbox(smtpCfg.From)
	if err != nil {
		return gerror.Wrap(err, "发件人格式不正确")
	}
	toAddr, err := parseAddressMailbox(to)
	if err != nil {
		return gerror.Wrap(err, "收件人格式不正确")
	}

	msg := buildEmailMessage(smtpCfg.From, to, subject, body)
	addr := fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port)
	auth := smtp.PlainAuth("", smtpCfg.User, smtpCfg.Pass, smtpCfg.Host)

	g.Log().Infof(ctx, "[notification:email] start host=%s port=%d to=%s", smtpCfg.Host, smtpCfg.Port, toAddr)

	if smtpCfg.Port == 465 {
		err = sendMailWithImplicitTLS(addr, smtpCfg.Host, auth, fromAddr, []string{toAddr}, msg)
	} else {
		err = smtp.SendMail(addr, auth, fromAddr, []string{toAddr}, msg)
	}
	if err != nil {
		g.Log().Errorf(ctx, "[notification:email] send failed host=%s port=%d to=%s err=%v", smtpCfg.Host, smtpCfg.Port, toAddr, err)
		return gerror.Wrap(err, "邮件发送失败")
	}
	g.Log().Infof(ctx, "[notification:email] send success host=%s port=%d to=%s", smtpCfg.Host, smtpCfg.Port, toAddr)
	return nil
}

func sendBySendGrid(ctx context.Context, sendgridCfg *SendGridConfig, to, subject, body string) error {
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{
			{
				"to": []map[string]string{
					{"email": strings.TrimSpace(to)},
				},
				"subject": subject,
			},
		},
		"from": map[string]string{
			"email": strings.TrimSpace(sendgridCfg.From),
		},
		"content": []map[string]string{
			{
				"type":  "text/plain",
				"value": body,
			},
		},
	}
	if strings.TrimSpace(sendgridCfg.FromName) != "" {
		payload["from"] = map[string]string{
			"email": strings.TrimSpace(sendgridCfg.From),
			"name":  strings.TrimSpace(sendgridCfg.FromName),
		}
	}
	return doSendGridRequest(ctx, sendgridCfg, payload)
}

func sendTemplateBySendGrid(ctx context.Context, sendgridCfg *SendGridConfig, to, subject, templateID string, params map[string]string) error {
	if strings.TrimSpace(templateID) == "" {
		return gerror.New("第三方模板ID不能为空")
	}
	dynamicData := map[string]interface{}{}
	for k, v := range params {
		dynamicData[k] = v
	}
	personalization := map[string]interface{}{
		"to": []map[string]string{
			{"email": strings.TrimSpace(to)},
		},
		"dynamic_template_data": dynamicData,
	}
	if strings.TrimSpace(subject) != "" {
		personalization["subject"] = subject
	}
	payload := map[string]interface{}{
		"personalizations": []map[string]interface{}{personalization},
		"from": map[string]string{
			"email": strings.TrimSpace(sendgridCfg.From),
		},
		"template_id": strings.TrimSpace(templateID),
	}
	if strings.TrimSpace(sendgridCfg.FromName) != "" {
		payload["from"] = map[string]string{
			"email": strings.TrimSpace(sendgridCfg.From),
			"name":  strings.TrimSpace(sendgridCfg.FromName),
		}
	}
	return doSendGridRequest(ctx, sendgridCfg, payload)
}

func doSendGridRequest(ctx context.Context, sendgridCfg *SendGridConfig, payload map[string]interface{}) error {
	if strings.TrimSpace(sendgridCfg.ApiKey) == "" {
		return gerror.New("SendGrid 配置缺少 API Key")
	}
	if strings.TrimSpace(sendgridCfg.From) == "" {
		return gerror.New("SendGrid 配置缺少发件人 from")
	}
	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return gerror.Wrap(err, "SendGrid 请求序列化失败")
	}
	req, err := http.NewRequest(http.MethodPost, "https://api.sendgrid.com/v3/mail/send", bytes.NewReader(bodyBytes))
	if err != nil {
		return gerror.Wrap(err, "创建 SendGrid 请求失败")
	}
	req = req.WithContext(ctx)
	req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(sendgridCfg.ApiKey))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		g.Log().Errorf(ctx, "[notification:sendgrid] request failed err=%v", err)
		return gerror.Wrap(err, "调用 SendGrid 失败")
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		g.Log().Errorf(ctx, "[notification:sendgrid] response failed status=%d body=%s", resp.StatusCode, string(respBody))
		return gerror.Newf("SendGrid 返回异常: status=%d body=%s", resp.StatusCode, string(respBody))
	}
	g.Log().Infof(ctx, "[notification:sendgrid] send success status=%d", resp.StatusCode)
	return nil
}

// SMSSender 短信发送占位实现。
type SMSSender struct{}

func (SMSSender) Send(ctx context.Context, to, body string) error {
	g.Log().Infof(ctx, "[notification:sms] to=%s body=%s", to, body)
	return nil
}

func parseAddressMailbox(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", gerror.New("地址不能为空")
	}
	addr, err := mail.ParseAddress(raw)
	if err == nil && addr != nil && strings.TrimSpace(addr.Address) != "" {
		return strings.TrimSpace(addr.Address), nil
	}
	// Fallback for plain mailbox strings.
	if strings.Contains(raw, "@") && !strings.ContainsAny(raw, " <>") {
		return raw, nil
	}
	return "", gerror.Wrap(err, "无法解析邮箱地址")
}

func buildEmailMessage(from, to, subject, body string) []byte {
	var b bytes.Buffer
	b.WriteString("To: " + to + "\r\n")
	b.WriteString("From: " + from + "\r\n")
	b.WriteString("Subject: " + mimeHeader(subject) + "\r\n")
	b.WriteString("MIME-Version: 1.0\r\n")
	b.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	b.WriteString("\r\n")
	b.WriteString(body)
	return b.Bytes()
}

func mimeHeader(subject string) string {
	// Keep ASCII-safe header to avoid introducing extra dependencies.
	// Non-ASCII subject can be added later with mime.QEncoding.
	return strings.ReplaceAll(subject, "\r", "")
}

func sendMailWithImplicitTLS(addr, host string, auth smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 10 * time.Second}, "tcp", addr, &tls.Config{
		ServerName: host,
		MinVersion: tls.VersionTLS12,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()

	if auth != nil {
		if ok, _ := client.Extension("AUTH"); ok {
			if err = client.Auth(auth); err != nil {
				return err
			}
		}
	}
	if err = client.Mail(from); err != nil {
		return err
	}
	for _, addrTo := range to {
		if err = client.Rcpt(addrTo); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(msg); err != nil {
		return err
	}
	if err = w.Close(); err != nil {
		return err
	}
	return client.Quit()
}
