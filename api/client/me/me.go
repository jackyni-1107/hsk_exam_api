package me

import (
	"context"

	v1 "exam/api/client/me/v1"
)

type IMe interface {
	Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error)
	MyExams(ctx context.Context, req *v1.ExamsReq) (res *v1.ExamsRes, err error)
}
