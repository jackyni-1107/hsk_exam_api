package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
	"exam/internal/consts"
	"exam/internal/middleware"
	membersvc "exam/internal/service/member"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

func (c *ControllerV1) MemberImport(ctx context.Context, req *v1.MemberImportReq) (res *v1.MemberImportRes, err error) {
	f := req.File
	if f == nil {
		f = g.RequestFromCtx(ctx).GetUploadFile("file")
	}
	if f == nil {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	rf, err := f.Open()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "无法读取上传文件")
	}
	defer rf.Close()

	creator := ""
	if d := middleware.GetCtxData(ctx); d != nil {
		creator = d.Username
	}
	out, err := membersvc.Member().MemberImport(ctx, rf, creator)
	if err != nil {
		return nil, err
	}
	return &v1.MemberImportRes{
		Total:   out.Total,
		Success: out.Success,
		Failed:  out.Failed,
		Errors:  out.Errors,
	}, nil
}
