package sysnotification

import (
	"context"
	"encoding/json"
	"strings"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	notifpkg "exam/internal/utility/notification"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (s *sSysNotification) Send(ctx context.Context, templateCode, channel, recipient, variables string) (bool, error) {
	g.Log().Infof(ctx, "[notification:send] start template=%s channel=%s recipient=%s", templateCode, channel, recipient)

	var tpl sysentity.SysNotificationTemplate
	err := dao.SysNotificationTemplate.Ctx(ctx).
		Where("code", templateCode).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&tpl)
	if err != nil {
		g.Log().Errorf(ctx, "[notification:send] load template failed template=%s err=%v", templateCode, err)
		return false, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if tpl.Id == 0 {
		g.Log().Warningf(ctx, "[notification:send] template not found template=%s", templateCode)
		return false, gerror.NewCode(consts.CodeTemplateNotFound)
	}

	vars := make(map[string]string)
	if variables != "" {
		if unmarshalErr := json.Unmarshal([]byte(variables), &vars); unmarshalErr != nil {
			g.Log().Warningf(ctx, "[notification:send] invalid variables json template=%s variables=%s err=%v", templateCode, variables, unmarshalErr)
		}
	}
	if tpl.TemplateType == 2 && strings.TrimSpace(tpl.ThirdPartyTemplateParams) != "" {
		defaultVars := make(map[string]string)
		if err := json.Unmarshal([]byte(tpl.ThirdPartyTemplateParams), &defaultVars); err != nil {
			g.Log().Warningf(ctx, "[notification:send] invalid third-party params template=%s params=%s err=%v", templateCode, tpl.ThirdPartyTemplateParams, err)
		} else {
			for k, v := range defaultVars {
				if _, exists := vars[k]; !exists {
					vars[k] = v
				}
			}
		}
	}
	rendered := notifpkg.RenderTemplate(tpl.Content, vars)
	subject := tpl.Name

	var sendErr error
	switch channel {
	case "email":
		if tpl.TemplateType == 2 {
			sendErr = (notifpkg.EmailSender{}).SendTemplate(ctx, recipient, subject, tpl.ThirdPartyTemplateId, vars)
		} else {
			sendErr = (notifpkg.EmailSender{}).Send(ctx, recipient, subject, rendered)
		}
	case "sms":
		sendErr = (notifpkg.SMSSender{}).Send(ctx, recipient, rendered)
	default:
		g.Log().Warningf(ctx, "[notification:send] unsupported channel template=%s channel=%s", templateCode, channel)
		return false, gerror.NewCode(consts.CodeUnsupportedChannel)
	}

	// log status aligns with admin page: 0=success, 1=failed.
	status := 0
	errMsg := ""
	if sendErr != nil {
		status = 1
		errMsg = sendErr.Error()
		g.Log().Errorf(ctx, "[notification:send] delivery failed template=%s channel=%s recipient=%s err=%v", templateCode, channel, recipient, sendErr)
	} else {
		g.Log().Infof(ctx, "[notification:send] delivery success template=%s channel=%s recipient=%s", templateCode, channel, recipient)
	}

	_, _ = dao.SysNotificationLog.Ctx(ctx).Insert(sysdo.SysNotificationLog{
		TemplateCode: templateCode,
		Channel:      channel,
		Recipient:    recipient,
		Content:      rendered,
		Status:       status,
		ErrorMsg:     errMsg,
	})

	if sendErr != nil {
		return false, gerror.WrapCode(consts.CodeInvalidParams, sendErr, "")
	}
	return true, nil
}
