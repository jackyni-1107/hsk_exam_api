package file

import (
	"context"

	v1 "exam/api/admin/file/v1"
	"exam/internal/middleware"
	sysfilesvc "exam/internal/service/sysfile"
	"exam/internal/utility"
)

func (c *ControllerV1) StorageConfigList(ctx context.Context, req *v1.StorageConfigListReq) (res *v1.StorageConfigListRes, err error) {
	list, err := sysfilesvc.SysFile().StorageConfigList(ctx)
	if err != nil {
		return nil, err
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
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
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
	id, err := sysfilesvc.SysFile().StorageConfigCreate(ctx, req.StorageType, req.Name, req.ConfigJson, creator, req.CleanupBeforeDays)
	if err != nil {
		return nil, err
	}
	return &v1.StorageConfigCreateRes{Id: id}, nil
}

func (c *ControllerV1) StorageConfigUpdate(ctx context.Context, req *v1.StorageConfigUpdateReq) (res *v1.StorageConfigUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysfilesvc.SysFile().StorageConfigUpdate(ctx, req.Id, req.Name, req.ConfigJson, updater, req.CleanupBeforeDays)
	if err != nil {
		return nil, err
	}
	return &v1.StorageConfigUpdateRes{}, nil
}

func (c *ControllerV1) StorageConfigDelete(ctx context.Context, req *v1.StorageConfigDeleteReq) (res *v1.StorageConfigDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysfilesvc.SysFile().StorageConfigDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.StorageConfigDeleteRes{}, nil
}

func (c *ControllerV1) StorageConfigSetActive(ctx context.Context, req *v1.StorageConfigSetActiveReq) (res *v1.StorageConfigSetActiveRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysfilesvc.SysFile().StorageConfigSetActive(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.StorageConfigSetActiveRes{}, nil
}
