package notification

import (
	"context"

	v1 "exam/api/admin/notification/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) TemplateList(ctx context.Context, req *v1.TemplateListReq) (res *v1.TemplateListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SysNotificationTemplate.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Code != "" {
		model = model.WhereLike("code", "%"+req.Code+"%")
	}
	if req.Channel != "" {
		model = model.Where("channel", req.Channel)
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysNotificationTemplate
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.TemplateItem, 0, len(list))
	for _, e := range list {
		item := &v1.TemplateItem{
			Id: int64(e.Id), Code: e.Code, Name: e.Name, Channel: e.Channel,
			Content: e.Content, Variables: e.Variables, Status: e.Status,
		}
		if e.CreateTime != nil {
			item.CreateTime = e.CreateTime.Format("Y-m-d H:i:s")
		}
		items = append(items, item)
	}
	return &v1.TemplateListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) TemplateCreate(ctx context.Context, req *v1.TemplateCreateReq) (res *v1.TemplateCreateRes, err error) {
	var exist sysentity.SysNotificationTemplate
	_ = dao.SysNotificationTemplate.Ctx(ctx).Where("code", req.Code).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&exist)
	if exist.Id > 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.template_exists")
	}
	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	id, err := dao.SysNotificationTemplate.Ctx(ctx).InsertAndGetId(sysdo.SysNotificationTemplate{
		Code: req.Code, Name: req.Name, Channel: req.Channel, Content: req.Content,
		Variables: req.Variables, Status: req.Status, Creator: creator, Updater: creator,
		DeleteFlag: consts.DeleteFlagNotDeleted,
	})
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.TemplateCreateRes{Id: id}, nil
}

func (c *ControllerV1) TemplateUpdate(ctx context.Context, req *v1.TemplateUpdateReq) (res *v1.TemplateUpdateRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	data := sysdo.SysNotificationTemplate{Updater: updater}
	if req.Name != "" {
		data.Name = req.Name
	}
	if req.Content != "" {
		data.Content = req.Content
	}
	if req.Variables != "" {
		data.Variables = req.Variables
	}
	data.Status = req.Status
	_, err = dao.SysNotificationTemplate.Ctx(ctx).Where("id", req.Id).Data(data).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.TemplateUpdateRes{}, nil
}

func (c *ControllerV1) TemplateDelete(ctx context.Context, req *v1.TemplateDeleteReq) (res *v1.TemplateDeleteRes, err error) {
	updater := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		updater = d.Username
	}
	_, err = dao.SysNotificationTemplate.Ctx(ctx).Where("id", req.Id).Data(sysdo.SysNotificationTemplate{
		DeleteFlag: consts.DeleteFlagDeleted, Updater: updater,
	}).Update()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return &v1.TemplateDeleteRes{}, nil
}
