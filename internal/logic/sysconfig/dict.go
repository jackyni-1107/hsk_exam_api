package sysconfig

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysconfig) DictTypeList(ctx context.Context, page, size int, dictType string) ([]sysentity.SysDictType, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	m := dao.SystemDictType.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if dictType != "" {
		m = m.WhereLike("dict_type", "%"+dictType+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysDictType
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysconfig) DictTypeCreate(ctx context.Context, dictName, dictType, remark, creator string, status int) (int64, error) {
	cnt, err := dao.SystemDictType.Ctx(ctx).
		Where("dict_type", dictType).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeDictTypeExists)
	}
	r, err := dao.SystemDictType.Ctx(ctx).Insert(sysdo.SysDictType{
		DictName:   dictName,
		DictType:   dictType,
		Remark:     remark,
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

func (s *sSysconfig) DictTypeUpdate(ctx context.Context, id int64, dictName, remark, updater string, status int) error {
	data := map[string]interface{}{
		"updater": updater,
		"status":  status,
	}
	if dictName != "" {
		data["dict_name"] = dictName
	}
	if remark != "" {
		data["remark"] = remark
	}
	_, err := dao.SystemDictType.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysconfig) DictTypeDelete(ctx context.Context, id int64, updater string) error {
	_, err := dao.SystemDictType.Ctx(ctx).Where("id", id).Data(sysdo.SysDictType{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysconfig) DictDataList(ctx context.Context, page, size int, dictType string) ([]sysentity.SysDictData, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 100
	}
	m := dao.SystemDictData.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if dictType != "" {
		m = m.Where("dict_type", dictType)
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysDictData
	err = m.Page(page, size).OrderAsc("sort").OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}

func (s *sSysconfig) DictDataCreate(ctx context.Context, dictType, dictLabel, dictValue, remark, creator string, sort, status int) (int64, error) {
	r, err := dao.SystemDictData.Ctx(ctx).Insert(sysdo.SysDictData{
		DictType:   dictType,
		DictLabel:  dictLabel,
		DictValue:  dictValue,
		Remark:     remark,
		Creator:    creator,
		Sort:       sort,
		Status:     status,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	id, _ := r.LastInsertId()
	return id, nil
}

func (s *sSysconfig) DictDataUpdate(ctx context.Context, id int64, dictLabel, dictValue, remark, updater string, sort, status int) error {
	data := map[string]interface{}{
		"updater": updater,
		"sort":    sort,
		"status":  status,
	}
	if dictLabel != "" {
		data["dict_label"] = dictLabel
	}
	if dictValue != "" {
		data["dict_value"] = dictValue
	}
	if remark != "" {
		data["remark"] = remark
	}
	_, err := dao.SystemDictData.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Data(data).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}

func (s *sSysconfig) DictDataDelete(ctx context.Context, id int64, updater string) error {
	_, err := dao.SystemDictData.Ctx(ctx).Where("id", id).Data(sysdo.SysDictData{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return nil
}
