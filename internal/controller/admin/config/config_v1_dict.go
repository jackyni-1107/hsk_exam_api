package config

import (
	"context"

	v1 "exam/api/admin/config/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) DictTypeList(ctx context.Context, req *v1.DictTypeListReq) (res *v1.DictTypeListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SystemDictType.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.DictType != "" {
		model = model.WhereLike("dict_type", "%"+req.DictType+"%")
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysDictType
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.DictTypeItem, 0, len(list))
	for _, e := range list {
		item := &v1.DictTypeItem{Id: int64(e.Id), DictName: e.DictName, DictType: e.DictType, Status: e.Status}
		if e.CreateTime != nil {
			item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.DictTypeListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) DictTypeCreate(ctx context.Context, req *v1.DictTypeCreateReq) (res *v1.DictTypeCreateRes, err error) {
	var exist sysentity.SysDictType
	_ = dao.SystemDictType.Ctx(ctx).Where("dict_type", req.DictType).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeDictTypeExists)
	}
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := dao.SystemDictType.Ctx(ctx).InsertAndGetId(sysdo.SysDictType{
		DictName: req.DictName, DictType: req.DictType, Status: req.Status, Remark: req.Remark,
		Creator: creator, Updater: creator, DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.DictTypeCreateRes{Id: id}, nil
}

func (c *ControllerV1) DictTypeUpdate(ctx context.Context, req *v1.DictTypeUpdateReq) (res *v1.DictTypeUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := sysdo.SysDictType{Updater: updater}
	if req.DictName != "" {
		data.DictName = req.DictName
	}
	data.Status = req.Status
	if req.Remark != "" {
		data.Remark = req.Remark
	}
	_, err = dao.SystemDictType.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.DictTypeUpdateRes{}, nil
}

func (c *ControllerV1) DictTypeDelete(ctx context.Context, req *v1.DictTypeDeleteReq) (res *v1.DictTypeDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SystemDictType.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysDictType{
		DeleteFlag: consts.DeleteFlagDeleted, Updater: updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.DictTypeDeleteRes{}, nil
}

func (c *ControllerV1) DictDataList(ctx context.Context, req *v1.DictDataListReq) (res *v1.DictDataListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 100
	}
	model := dao.SystemDictData.Ctx(ctx).Where("dict_type", req.DictType).Where("delete_flag", consts.DeleteFlagNotDeleted)
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysDictData
	err = model.Page(req.Page, req.Size).OrderAsc("sort").OrderAsc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.DictDataItem, 0, len(list))
	for _, e := range list {
		item := &v1.DictDataItem{Id: int64(e.Id), DictType: e.DictType, DictLabel: e.DictLabel, DictValue: e.DictValue, Sort: e.Sort, Status: e.Status}
		if e.CreateTime != nil {
			item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.DictDataListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) DictDataCreate(ctx context.Context, req *v1.DictDataCreateReq) (res *v1.DictDataCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := dao.SystemDictData.Ctx(ctx).InsertAndGetId(sysdo.SysDictData{
		DictType: req.DictType, DictLabel: req.DictLabel, DictValue: req.DictValue,
		Sort: req.Sort, Status: req.Status, Remark: req.Remark,
		Creator: creator, Updater: creator, DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.DictDataCreateRes{Id: id}, nil
}

func (c *ControllerV1) DictDataUpdate(ctx context.Context, req *v1.DictDataUpdateReq) (res *v1.DictDataUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := sysdo.SysDictData{Updater: updater}
	if req.DictLabel != "" {
		data.DictLabel = req.DictLabel
	}
	if req.DictValue != "" {
		data.DictValue = req.DictValue
	}
	data.Sort = req.Sort
	data.Status = req.Status
	if req.Remark != "" {
		data.Remark = req.Remark
	}
	_, err = dao.SystemDictData.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.DictDataUpdateRes{}, nil
}

func (c *ControllerV1) DictDataDelete(ctx context.Context, req *v1.DictDataDeleteReq) (res *v1.DictDataDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SystemDictData.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysDictData{
		DeleteFlag: consts.DeleteFlagDeleted, Updater: updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.DictDataDeleteRes{}, nil
}
