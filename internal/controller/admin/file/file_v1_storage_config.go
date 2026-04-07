package file

import (
	"context"

	v1 "exam/api/admin/file/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) StorageConfigList(ctx context.Context, req *v1.StorageConfigListReq) (res *v1.StorageConfigListRes, err error) {
	var list []sysentity.SysFileStorageConfig
	err = dao.SysFileStorageConfig.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.StorageConfigItem, 0, len(list))
	for _, e := range list {
		item := &v1.StorageConfigItem{
			Id:                int64(e.Id),
			StorageType:       e.StorageType,
			Name:              e.Name,
			IsActive:          e.IsActive,
			ConfigJson:        e.ConfigJson,
			CleanupBeforeDays: e.CleanupBeforeDays,
		}
		if e.CreateTime != nil {
			item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.StorageConfigListRes{List: items}, nil
}

func (c *ControllerV1) StorageConfigCreate(ctx context.Context, req *v1.StorageConfigCreateReq) (res *v1.StorageConfigCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	cleanupDays := req.CleanupBeforeDays
	if cleanupDays <= 0 {
		cleanupDays = 30
	}
	id, err := dao.SysFileStorageConfig.Ctx(ctx).InsertAndGetId(sysdo.SysFileStorageConfig{
		StorageType:       req.StorageType,
		Name:              req.Name,
		IsActive:          0,
		ConfigJson:        req.ConfigJson,
		CleanupBeforeDays: cleanupDays,
		Creator:           creator,
		Updater:           creator,
		DeleteFlag:        consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.StorageConfigCreateRes{Id: id}, nil
}

func (c *ControllerV1) StorageConfigUpdate(ctx context.Context, req *v1.StorageConfigUpdateReq) (res *v1.StorageConfigUpdateRes, err error) {
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
	if req.CleanupBeforeDays > 0 {
		data["cleanup_before_days"] = req.CleanupBeforeDays
	}
	_, err = dao.SysFileStorageConfig.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.StorageConfigUpdateRes{}, nil
}

func (c *ControllerV1) StorageConfigDelete(ctx context.Context, req *v1.StorageConfigDeleteReq) (res *v1.StorageConfigDeleteRes, err error) {
	var e sysentity.SysFileStorageConfig
	err = dao.SysFileStorageConfig.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.config_not_found")
	}
	if e.IsActive == 1 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.cannot_delete_active_config")
	}
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SysFileStorageConfig.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysFileStorageConfig{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.StorageConfigDeleteRes{}, nil
}

func (c *ControllerV1) StorageConfigSetActive(ctx context.Context, req *v1.StorageConfigSetActiveReq) (res *v1.StorageConfigSetActiveRes, err error) {
	var e sysentity.SysFileStorageConfig
	err = dao.SysFileStorageConfig.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil || e.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.config_not_found")
	}
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	// 先全部设为未启用
	_, _ = dao.SysFileStorageConfig.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(sysdo.SysFileStorageConfig{IsActive: 0, Updater: updater}).
		Update()
	// 再启用当前
	_, err = dao.SysFileStorageConfig.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysFileStorageConfig{
		IsActive: 1,
		Updater:  updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.StorageConfigSetActiveRes{}, nil
}
