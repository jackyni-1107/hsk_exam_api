package file

import (
	"context"

	"exam/api/admin/file/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/do"
	"exam/internal/model/entity"
	"exam/internal/storage"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	model := dao.SysFileStorage.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if req.Filename != "" {
		model = model.WhereLike("filename", "%"+req.Filename+"%")
	}
	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []entity.SysFileStorage
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	items := make([]*v1.FileItem, 0, len(list))
	for _, e := range list {
		item := &v1.FileItem{
			Id: int64(e.Id), Filename: e.Filename, Path: e.Path, Size: e.Size,
			MimeType: e.MimeType, IsPrivate: e.IsPrivate,
		}
		if e.CreateTime != nil {
			item.CreateTime = e.CreateTime.Format("Y-m-d H:i:s")
		}
		items = append(items, item)
	}
	return &v1.FileListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error) {
	var f entity.SysFileStorage
	err = dao.SysFileStorage.Ctx(ctx).Where("id", req.Id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&f)
	if err != nil || f.Id == 0 {
		return nil, gerror.NewCode(consts.CodeInvalidParams, "err.file_not_found")
	}
	adapter := storage.NewAdapter()
	_ = adapter.Delete(ctx, f.Bucket, f.Path)
	_, _ = dao.SysFileStorage.Ctx(ctx).Where("id", req.Id).Data(do.SysFileStorage{DeleteFlag: consts.DeleteFlagDeleted}).Update()
	return &v1.FileDeleteRes{}, nil
}
