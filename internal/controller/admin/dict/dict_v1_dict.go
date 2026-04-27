package dict

import (
	"context"

	v1 "exam/api/admin/dict/v1"
	"exam/internal/middleware"
	sysdictsvc "exam/internal/service/sysdict"
	"exam/internal/utility"
)

func (c *ControllerV1) DictTypeList(ctx context.Context, req *v1.DictTypeListReq) (res *v1.DictTypeListRes, err error) {
	list, total, err := sysdictsvc.SysDict().DictTypeList(ctx, req.Page, req.Size, req.DictType)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.DictTypeItem, 0, len(list))
	for _, e := range list {
		item := &v1.DictTypeItem{Id: int64(e.Id), DictName: e.DictName, DictType: e.DictType, Status: e.Status}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.DictTypeListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) DictTypeCreate(ctx context.Context, req *v1.DictTypeCreateReq) (res *v1.DictTypeCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := sysdictsvc.SysDict().DictTypeCreate(ctx, req.DictName, req.DictType, req.Remark, creator, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.DictTypeCreateRes{Id: id}, nil
}

func (c *ControllerV1) DictTypeUpdate(ctx context.Context, req *v1.DictTypeUpdateReq) (res *v1.DictTypeUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysdictsvc.SysDict().DictTypeUpdate(ctx, req.Id, req.DictName, req.Remark, updater, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.DictTypeUpdateRes{}, nil
}

func (c *ControllerV1) DictTypeDelete(ctx context.Context, req *v1.DictTypeDeleteReq) (res *v1.DictTypeDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysdictsvc.SysDict().DictTypeDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.DictTypeDeleteRes{}, nil
}

func (c *ControllerV1) DictDataList(ctx context.Context, req *v1.DictDataListReq) (res *v1.DictDataListRes, err error) {
	list, total, err := sysdictsvc.SysDict().DictDataList(ctx, req.Page, req.Size, req.DictType)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.DictDataItem, 0, len(list))
	for _, e := range list {
		item := &v1.DictDataItem{Id: int64(e.Id), DictType: e.DictType, DictLabel: e.DictLabel, DictValue: e.DictValue, Sort: e.Sort, Status: e.Status}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
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
	id, err := sysdictsvc.SysDict().DictDataCreate(ctx, req.DictType, req.DictLabel, req.DictValue, req.Remark, creator, req.Sort, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.DictDataCreateRes{Id: id}, nil
}

func (c *ControllerV1) DictDataUpdate(ctx context.Context, req *v1.DictDataUpdateReq) (res *v1.DictDataUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysdictsvc.SysDict().DictDataUpdate(ctx, req.Id, req.DictLabel, req.DictValue, req.Remark, updater, req.Sort, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.DictDataUpdateRes{}, nil
}

func (c *ControllerV1) DictDataDelete(ctx context.Context, req *v1.DictDataDeleteReq) (res *v1.DictDataDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = sysdictsvc.SysDict().DictDataDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.DictDataDeleteRes{}, nil
}
