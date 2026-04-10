package notification

import (
	"context"
	"strings"

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

// EmailSender SMTP 发送占位实现（按渠道配置扩展）。
type EmailSender struct{}

func (EmailSender) Send(ctx context.Context, to, body string) error {
	g.Log().Infof(ctx, "[notification:email] to=%s body=%s", to, body)
	return nil
}

// SMSSender 短信发送占位实现。
type SMSSender struct{}

func (SMSSender) Send(ctx context.Context, to, body string) error {
	g.Log().Infof(ctx, "[notification:sms] to=%s body=%s", to, body)
	return nil
}
