package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) ChannelConfigList(ctx context.Context, req *v1.ChannelConfigListReq) (res *v1.ChannelConfigListRes, err error) {
	model := dao.SysNotificationChannelConfig.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Channel != "" {
		model = model.Where("channel", req.Channel)
	}
	var list []sysentity.SysNotificationChannelConfig
	err = model.OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.ChannelConfigItem, 0, len(list))
	for _, e := range list {
		item := &v1.ChannelConfigItem{
			Id:         int64(e.Id),
			Channel:    e.Channel,
			Provider:   e.Provider,
			Name:       e.Name,
			IsActive:   e.IsActive,
			ConfigJson: e.ConfigJson,
		}
		if e.CreateTime != nil {
			item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.ChannelConfigListRes{List: items}, nil
}

func (c *ControllerV1) ChannelConfigCreate(ctx context.Context, req *v1.ChannelConfigCreateReq) (res *v1.ChannelConfigCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	// 校验 provider 与 channel 匹配
	if req.Channel == "email" && req.Provider != "smtp" {
		return nil, gerror.NewCode(consts.CodeEmailMustUseSmtp)
	}
	if req.Channel == "sms" && req.Provider != "aliyun" && req.Provider != "tencent" {
		return nil, gerror.NewCode(consts.CodeSmsMustUseAliyunOrTencent)
	}
	id, err := dao.SysNotificationChannelConfig.Ctx(ctx).InsertAndGetId(sysdo.SysNotificationChannelConfig{
		Channel:    req.Channel,
		Provider:   req.Provider,
		Name:       req.Name,
		IsActive:   0,
		ConfigJson: req.ConfigJson,
		Creator:    creator,
		Updater:    creator,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ChannelConfigCreateRes{Id: id}, nil
}

func (c *ControllerV1) ChannelConfigUpdate(ctx context.Context, req *v1.ChannelConfigUpdateReq) (res *v1.ChannelConfigUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := map[string]interface{}{"updater": updater}
	if req.Name != "" {
		data["name"] = req.Name
	}
	if req.ConfigJson != "" {
		data["config_json"] = req.ConfigJson
	}
	_, err = dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ChannelConfigUpdateRes{}, nil
}

func (c *ControllerV1) ChannelConfigDelete(ctx context.Context, req *v1.ChannelConfigDeleteReq) (res *v1.ChannelConfigDeleteRes, err error) {
	var e sysentity.SysNotificationChannelConfig
	err = dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeConfigNotFound)
	}
	if e.IsActive == 1 {
		return nil, gerror.NewCode(consts.CodeCannotDeleteActiveConfig)
	}
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysNotificationChannelConfig{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ChannelConfigDeleteRes{}, nil
}

func (c *ControllerV1) ChannelConfigSetActive(ctx context.Context, req *v1.ChannelConfigSetActiveReq) (res *v1.ChannelConfigSetActiveRes, err error) {
	var e sysentity.SysNotificationChannelConfig
	err = dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeConfigNotFound)
	}
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	// 同渠道先全部设为未启用
	_, _ = dao.SysNotificationChannelConfig.Ctx(ctx).
		Where("channel", e.Channel).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(sysdo.SysNotificationChannelConfig{IsActive: 0, Updater: updater}).
		Update()
	// 再启用当前
	_, err = dao.SysNotificationChannelConfig.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysNotificationChannelConfig{
		IsActive: 1,
		Updater:  updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ChannelConfigSetActiveRes{}, nil
}
