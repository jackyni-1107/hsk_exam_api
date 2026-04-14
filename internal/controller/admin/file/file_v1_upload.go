package file

import (
	"context"

	v1 "exam/api/admin/file/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	SysFilesvc "exam/internal/service/SysFile"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) Upload(ctx context.Context, req *v1.FileUploadReq) (res *v1.FileUploadRes, err error) {
	f := req.File
	if f == nil {
		f = g.RequestFromCtx(ctx).GetUploadFile("file")
	}
	if f == nil {
		return nil, gerror.NewCode(consts.CodeFileRequired)
	}
	rf, err := f.Open()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	defer rf.Close()

	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	ct := f.Header.Get("Content-Type")
	if ct == "" {
		ct = "application/octet-stream"
	}
	id, objectPath, displayName, err := SysFilesvc.SysFile().FileUpload(ctx, f.Filename, f.Size, ct, rf, req.IsPrivate, creator)
	if err != nil {
		return nil, err
	}
	return &v1.FileUploadRes{
		Id:       id,
		Path:     objectPath,
		Filename: displayName,
		Size:     f.Size,
		MimeType: ct,
	}, nil
}
