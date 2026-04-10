package sysfile

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/utility/storage"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysfile) FileList(ctx context.Context, page, size int, filename string) ([]sysentity.SysFileStorage, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	m := dao.SysFileStorage.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if filename != "" {
		m = m.WhereLike("filename", "%"+filename+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysFileStorage
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysfile) FileDelete(ctx context.Context, id int64) error {
	var e sysentity.SysFileStorage
	err := dao.SysFileStorage.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return gerror.NewCode(consts.CodeFileNotFound)
	}
	adapter := storage.NewAdapter()
	_ = adapter.Delete(ctx, e.Bucket, e.Path)
	_, err = dao.SysFileStorage.Ctx(ctx).Where("id", id).Data(sysdo.SysFileStorage{
		DeleteFlag: consts.DeleteFlagDeleted,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysfile) StorageConfigList(ctx context.Context) ([]sysentity.SysFileStorageConfig, error) {
	var list []sysentity.SysFileStorageConfig
	err := dao.SysFileStorageConfig.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, nil
}

func (s *sSysfile) StorageConfigCreate(ctx context.Context, storageType, name, configJson, creator string, cleanupBeforeDays int) (int64, error) {
	if cleanupBeforeDays <= 0 {
		cleanupBeforeDays = 30
	}
	r, err := dao.SysFileStorageConfig.Ctx(ctx).Insert(sysdo.SysFileStorageConfig{
		StorageType:       storageType,
		Name:              name,
		ConfigJson:        configJson,
		Creator:           creator,
		CleanupBeforeDays: cleanupBeforeDays,
		IsActive:          0,
		DeleteFlag:        consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	return id, nil
}

func (s *sSysfile) StorageConfigUpdate(ctx context.Context, id int64, name, configJson, updater string, cleanupBeforeDays int) error {
	data := map[string]interface{}{
		"updater": updater,
	}
	if name != "" {
		data["name"] = name
	}
	if configJson != "" {
		data["config_json"] = configJson
	}
	if cleanupBeforeDays > 0 {
		data["cleanup_before_days"] = cleanupBeforeDays
	}
	_, err := dao.SysFileStorageConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysfile) StorageConfigDelete(ctx context.Context, id int64, updater string) error {
	var e sysentity.SysFileStorageConfig
	err := dao.SysFileStorageConfig.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&e)
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if e.Id == 0 {
		return gerror.NewCode(consts.CodeConfigNotFound)
	}
	if e.IsActive == 1 {
		return gerror.NewCode(consts.CodeCannotDeleteActiveConfig)
	}
	_, err = dao.SysFileStorageConfig.Ctx(ctx).Where("id", id).Data(sysdo.SysFileStorageConfig{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysfile) StorageConfigSetActive(ctx context.Context, id int64, updater string) error {
	_, err := dao.SysFileStorageConfig.Ctx(ctx).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(map[string]interface{}{"is_active": 0, "updater": updater}).
		Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	_, err = dao.SysFileStorageConfig.Ctx(ctx).
		Where("id", id).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(map[string]interface{}{"is_active": 1, "updater": updater}).
		Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}
