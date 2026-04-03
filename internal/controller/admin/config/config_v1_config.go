package config

import (
	"context"

	"exam/api/admin/config/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	"exam/internal/model/do"
	"exam/internal/model/entity"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) ConfigList(ctx context.Context, req *v1.ConfigListReq) (res *v1.ConfigListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SystemConfig.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Group != "" {
		model = model.Where("group_name", req.Group)
	}
	if req.Key != "" {
		model = model.WhereLike("config_key", "%"+req.Key+"%")
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []entity.SystemConfig
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.ConfigItem, 0, len(list))
	for _, e := range list {
		item := &v1.ConfigItem{
			Id: int64(e.Id), ConfigKey: e.ConfigKey, ConfigValue: e.ConfigValue,
			ConfigType: e.ConfigType, GroupName: e.GroupName, Remark: e.Remark,
		}
		if e.CreateTime != nil {
			item.CreateTime = e.CreateTime.Format("Y-m-d H:i:s")
		}
		items = append(items, item)
	}
	return &v1.ConfigListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) ConfigCreate(ctx context.Context, req *v1.ConfigCreateReq) (res *v1.ConfigCreateRes, err error) {
	var exist entity.SystemConfig
	_ = dao.SystemConfig.Ctx(ctx).Where("config_key", req.ConfigKey).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.config_exists")
	}
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	configType := req.ConfigType
	if configType == "" {
		configType = "string"
	}
	groupName := req.GroupName
	if groupName == "" {
		groupName = "default"
	}
	id, err := dao.SystemConfig.Ctx(ctx).InsertAndGetId(do.SystemConfig{
		ConfigKey: req.ConfigKey, ConfigValue: req.ConfigValue, ConfigType: configType,
		GroupName: groupName, Remark: req.Remark, Creator: creator, Updater: creator,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ConfigCreateRes{Id: id}, nil
}

func (c *ControllerV1) ConfigUpdate(ctx context.Context, req *v1.ConfigUpdateReq) (res *v1.ConfigUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := do.SystemConfig{Updater: updater}
	if req.ConfigValue != "" {
		data.ConfigValue = req.ConfigValue
	}
	if req.Remark != "" {
		data.Remark = req.Remark
	}
	_, err = dao.SystemConfig.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ConfigUpdateRes{}, nil
}

func (c *ControllerV1) ConfigDelete(ctx context.Context, req *v1.ConfigDeleteReq) (res *v1.ConfigDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SystemConfig.Ctx(ctx).Where("id", req.Id).Data(do.SystemConfig{
		DeleteFlag: consts.DeleteFlagDeleted, Updater: updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.ConfigDeleteRes{}, nil
}

func (c *ControllerV1) ConfigGet(ctx context.Context, req *v1.ConfigGetReq) (res *v1.ConfigGetRes, err error) {
	var e entity.SystemConfig
	err = dao.SystemConfig.Ctx(ctx).Where("config_key", req.Key).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.config_not_found")
	}
	return &v1.ConfigGetRes{Value: e.ConfigValue}, nil
}
