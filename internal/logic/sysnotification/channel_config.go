package sysnotification

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysnotification) ChannelConfigList(ctx context.Context, channel string) ([]sysentity.SysNotificationChannelConfig, error) {
	m := dao.SysNotificationChannelConfig.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if channel != "" {
		m = m.Where("channel", channel)
	}
	var list []sysentity.SysNotificationChannelConfig
	err := m.OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, nil
}

func (s *sSysnotification) ChannelConfigCreate(ctx context.Context, channel, provider, name, configJson, creator string) (int64, error) {
	switch channel {
	case "email":
		if provider != "smtp" {
			return 0, gerror.NewCode(consts.CodeEmailMustUseSmtp)
		}
	case "sms":
		if provider != "aliyun" && provider != "tencent" {
			return 0, gerror.NewCode(consts.CodeSmsMustUseAliyunOrTencent)
		}
	default:
		return 0, gerror.NewCode(consts.CodeUnsupportedChannel)
	}
	r, err := dao.SysNotificationChannelConfig.Ctx(ctx).Insert(sysdo.SysNotificationChannelConfig{
		Channel:    channel,
		Provider:   provider,
		Name:       name,
		ConfigJson: configJson,
		Creator:    creator,
		IsActive:   0,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	return id, nil
}

func (s *sSysnotification) ChannelConfigUpdate(ctx context.Context, id int64, name, configJson, updater string) error {
	data := map[string]interface{}{
		"updater": updater,
	}
	if name != "" {
		data["name"] = name
	}
	if configJson != "" {
		data["config_json"] = configJson
	}
	_, err := dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysnotification) ChannelConfigDelete(ctx context.Context, id int64, updater string) error {
	var e sysentity.SysNotificationChannelConfig
	err := dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return gerror.NewCode(consts.CodeConfigNotFound)
	}
	if e.IsActive == 1 {
		return gerror.NewCode(consts.CodeCannotDeleteActiveConfig)
	}
	_, err = dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", id).Data(sysdo.SysNotificationChannelConfig{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysnotification) ChannelConfigSetActive(ctx context.Context, id int64, updater string) error {
	var e sysentity.SysNotificationChannelConfig
	err := dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return gerror.NewCode(consts.CodeConfigNotFound)
	}
	_, err = dao.SysNotificationChannelConfig.Ctx(ctx).
		Where("channel", e.Channel).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(map[string]interface{}{"is_active": 0, "updater": updater}).
		Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	_, err = dao.SysNotificationChannelConfig.Ctx(ctx).
		Where("id", id).
		Data(map[string]interface{}{"is_active": 1, "updater": updater}).
		Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}
