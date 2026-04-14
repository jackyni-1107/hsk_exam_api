package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
	"exam/internal/middleware"
	notisvc "exam/internal/service/SysNotification"
	"exam/internal/utility"
)

func (c *ControllerV1) TemplateList(ctx context.Context, req *v1.TemplateListReq) (res *v1.TemplateListRes, err error) {
	list, total, err := notisvc.SysNotification().TemplateList(ctx, req.Page, req.Size, req.Code, req.Channel)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.TemplateItem, 0, len(list))
	for _, e := range list {
		item := &v1.TemplateItem{
			Id: int64(e.Id), Code: e.Code, Name: e.Name, Channel: e.Channel,
			Content: e.Content, Variables: e.Variables, Status: e.Status,
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.TemplateListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) TemplateCreate(ctx context.Context, req *v1.TemplateCreateReq) (res *v1.TemplateCreateRes, err error) {
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := notisvc.SysNotification().TemplateCreate(ctx, req.Code, req.Name, req.Channel, req.Content, req.Variables, creator, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.TemplateCreateRes{Id: id}, nil
}

func (c *ControllerV1) TemplateUpdate(ctx context.Context, req *v1.TemplateUpdateReq) (res *v1.TemplateUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = notisvc.SysNotification().TemplateUpdate(ctx, req.Id, req.Name, req.Content, req.Variables, updater, req.Status)
	if err != nil {
		return nil, err
	}
	return &v1.TemplateUpdateRes{}, nil
}

func (c *ControllerV1) TemplateDelete(ctx context.Context, req *v1.TemplateDeleteReq) (res *v1.TemplateDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	err = notisvc.SysNotification().TemplateDelete(ctx, req.Id, updater)
	if err != nil {
		return nil, err
	}
	return &v1.TemplateDeleteRes{}, nil
}
