package sysnotification

import (
	"context"
	"encoding/json"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	notifpkg "exam/internal/utility/notification"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysnotification) Send(ctx context.Context, templateCode, channel, recipient, variables string) (bool, error) {
	var tpl sysentity.SysNotificationTemplate
	err := dao.SysNotificationTemplate.Ctx(ctx).
		Where("code", templateCode).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&tpl)
	if err != nil {
		return false, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if tpl.Id == 0 {
		return false, gerror.NewCode(consts.CodeTemplateNotFound)
	}

	vars := make(map[string]string)
	if variables != "" {
		_ = json.Unmarshal([]byte(variables), &vars)
	}
	rendered := notifpkg.RenderTemplate(tpl.Content, vars)

	var sendErr error
	switch channel {
	case "email":
		sendErr = (notifpkg.EmailSender{}).Send(ctx, recipient, rendered)
	case "sms":
		sendErr = (notifpkg.SMSSender{}).Send(ctx, recipient, rendered)
	default:
		return false, gerror.NewCode(consts.CodeUnsupportedChannel)
	}

	status := 1
	errMsg := ""
	if sendErr != nil {
		status = 2
		errMsg = sendErr.Error()
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
