package sysnotification

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"exam/internal/auditutil"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysNotification) TemplateList(ctx context.Context, page, size int, code, channel string) ([]sysentity.SysNotificationTemplate, int, error) {
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

func (s *sSysNotification) TemplateCreate(ctx context.Context, code, name, channel, content, variables, creator string, status int) (int64, error) {
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
	var after sysentity.SysNotificationTemplate
	if err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SysNotificationTemplate.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sSysNotification) TemplateUpdate(ctx context.Context, id int64, name, content, variables, updater string, status int) error {
	var before sysentity.SysNotificationTemplate
	if err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeTemplateNotFound)
	}
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
	data["variables"] = variables
	_, err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var after sysentity.SysNotificationTemplate
	if err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysNotificationTemplate.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysNotification) TemplateDelete(ctx context.Context, id int64, updater string) error {
	var before sysentity.SysNotificationTemplate
	if err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeTemplateNotFound)
	}
	_, err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Data(sysdo.SysNotificationTemplate{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var after sysentity.SysNotificationTemplate
	if err := dao.SysNotificationTemplate.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysNotificationTemplate.Table(), id, &before, &after)
	}
	return nil
}
