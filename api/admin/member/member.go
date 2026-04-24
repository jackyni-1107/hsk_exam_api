package member

import (
	"context"

	v1 "exam/api/admin/member/v1"
)

type IMember interface {
	MemberList(ctx context.Context, req *v1.MemberListReq) (res *v1.MemberListRes, err error)
	MemberCreate(ctx context.Context, req *v1.MemberCreateReq) (res *v1.MemberCreateRes, err error)
	MemberUpdate(ctx context.Context, req *v1.MemberUpdateReq) (res *v1.MemberUpdateRes, err error)
	MemberDelete(ctx context.Context, req *v1.MemberDeleteReq) (res *v1.MemberDeleteRes, err error)
	MemberImport(ctx context.Context, req *v1.MemberImportReq) (res *v1.MemberImportRes, err error)
}
