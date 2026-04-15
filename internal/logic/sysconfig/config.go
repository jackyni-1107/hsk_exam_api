package sysconfig

import (
	"context"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysConfig) ConfigList(ctx context.Context, page, size int, group, key string) ([]sysentity.SysConfig, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	m := dao.SystemConfig.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if group != "" {
		m = m.Where("group_name", group)
	}
	if key != "" {
		m = m.WhereLike("config_key", "%"+key+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysConfig
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysConfig) ConfigCreate(ctx context.Context, configKey, configValue, configType, groupName, remark, creator string) (int64, error) {
	cnt, err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", configKey).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeConfigExists)
	}
	if configType == "" {
		configType = "string"
	}
	if groupName == "" {
		groupName = "default"
	}
	r, err := dao.SystemConfig.Ctx(ctx).Insert(sysdo.SysConfig{
		ConfigKey:   configKey,
		ConfigValue: configValue,
		ConfigType:  configType,
		GroupName:   groupName,
		Remark:      remark,
		Creator:     creator,
		DeleteFlag:  consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	var after sysentity.SysConfig
	if err := dao.SystemConfig.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SystemConfig.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sSysConfig) ConfigUpdate(ctx context.Context, id int64, configValue, remark, updater string) error {
	var before sysentity.SysConfig
	if err := dao.SystemConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeConfigNotFound)
	}
	data := map[string]interface{}{
		"updater": updater,
	}
	if configValue != "" {
		data["config_value"] = configValue
	}
	if remark != "" {
		data["remark"] = remark
	}
	_, err := dao.SystemConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var after sysentity.SysConfig
	if err := dao.SystemConfig.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemConfig.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysConfig) ConfigDelete(ctx context.Context, id int64, updater string) error {
	var before sysentity.SysConfig
	if err := dao.SystemConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeConfigNotFound)
	}
	_, err := dao.SystemConfig.Ctx(ctx).Where("id", id).Data(sysdo.SysConfig{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var after sysentity.SysConfig
	if err := dao.SystemConfig.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemConfig.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysConfig) ConfigGet(ctx context.Context, key string) (string, error) {
	var e sysentity.SysConfig
	err := dao.SystemConfig.Ctx(ctx).
		Where("config_key", key).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&e)
	if err != nil {
		return "", gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return "", gerror.NewCode(consts.CodeConfigNotFound)
	}
	return e.ConfigValue, nil
}
