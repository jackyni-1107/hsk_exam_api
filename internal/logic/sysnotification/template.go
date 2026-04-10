package sysnotification

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysnotification) TemplateList(ctx context.Context, page, size int, code, channel string) ([]sysentity.SysNotificationTemplate, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	m := dao.SysNotificationTemplate.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if code != "" {
		m = m.WhereLike("code", "%"+code+"%")
	}
	if channel != "" {
		m = m.Where("channel", channel)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysNotificationTemplate
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysnotification) TemplateCreate(ctx context.Context, code, name, channel, content, variables, creator string, status int) (int64, error) {
	cnt, err := dao.SysNotificationTemplate.Ctx(ctx).
		Where("code", code).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeTemplateExists)
	}
	r, err := dao.SysNotificationTemplate.Ctx(ctx).Insert(sysdo.SysNotificationTemplate{
		Code:       code,
		Name:       name,
		Channel:    channel,
		Content:    content,
		Variables:  variables,
		Creator:    creator,
		Status:     status,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	return id, nil
}

func (s *sSysnotification) TemplateUpdate(ctx context.Context, id int64, name, content, variables, updater string, status int) error {
	data := map[string]interface{}{
		"updater": updater,
		"status":  status,
	}
	if name != "" {
		data["name"] = name
	}
	if content != "" {
		data["content"] = content
	}
	if variables != "" {
		data["variables"] = variables
	}
	_, err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysnotification) TemplateDelete(ctx context.Context, id int64, updater string) error {
	_, err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Data(sysdo.SysNotificationTemplate{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}
