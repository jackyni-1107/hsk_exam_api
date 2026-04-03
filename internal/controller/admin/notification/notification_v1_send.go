package notification

import (
	"context"
	"encoding/json"

	"exam/api/admin/notification/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	dosys "exam/internal/model/do/sys"
	"exam/internal/model/entity"
	notifpkg "exam/internal/notification"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) Send(ctx context.Context, req *v1.NotificationSendReq) (res *v1.NotificationSendRes, err error) {
	var tpl entity.SysNotificationTemplate
	err = dao.SysNotificationTemplate.Ctx(ctx).
		Where("code", req.TemplateCode).
		Where("channel", req.Channel).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&tpl)
	if err != nil || tpl.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.template_not_found")
	}
	vars := map[string]string{}
	if req.Variables != "" {
		_ = json.Unmarshal([]byte(req.Variables), &vars)
	}
	body := notifpkg.RenderTemplate(tpl.Content, vars)
	status := 1
	errMsg := ""
	switch req.Channel {
	case "email":
		err = notifpkg.EmailSender{}.Send(ctx, req.Recipient, body)
	case "sms":
		err = notifpkg.SMSSender{}.Send(ctx, req.Recipient, body)
	default:
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.unsupported_channel")
	}
	if err != nil {
		status = 2
		errMsg = err.Error()
	}
	_, _ = dao.SysNotificationLog.Ctx(ctx).Insert(dosys.SysNotificationLog{
		TemplateCode: req.TemplateCode,
		Channel:      req.Channel,
		Recipient:    req.Recipient,
		Content:      body,
		Status:       status,
		ErrorMsg:     errMsg,
	})
	if status != 1 {
		return &v1.NotificationSendRes{Ok: false}, nil
	}
	return &v1.NotificationSendRes{Ok: true}, nil
}
