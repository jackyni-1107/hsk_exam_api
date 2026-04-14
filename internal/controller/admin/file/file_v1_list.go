package file

import (
	"context"

	v1 "exam/api/admin/file/v1"
	sysfilesvc "exam/internal/service/sysfile"
	"exam/internal/utility"
)

func (c *ControllerV1) List(ctx context.Context, req *v1.FileListReq) (res *v1.FileListRes, err error) {
	list, total, err := sysfilesvc.SysFile().FileList(ctx, req.Page, req.Size, req.Filename)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.FileItem, 0, len(list))
	for _, e := range list {
		item := &v1.FileItem{
			Id: int64(e.Id), Filename: e.Filename, Path: e.Path, Size: e.Size,
			MimeType: e.MimeType, IsPrivate: e.IsPrivate,
		}
		if e.CreateTime != nil {
			item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		}
		items = append(items, item)
	}
	return &v1.FileListRes{List: items, Total: total}, nil
}

func (c *ControllerV1) Delete(ctx context.Context, req *v1.FileDeleteReq) (res *v1.FileDeleteRes, err error) {
	err = sysfilesvc.SysFile().FileDelete(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.FileDeleteRes{}, nil
}
