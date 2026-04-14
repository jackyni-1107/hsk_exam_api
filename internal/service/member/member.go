// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package member

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	IMember interface {
		MemberList(ctx context.Context, page int, size int, username string, status int) ([]sysentity.SysMember, int, error)
		MemberCreate(ctx context.Context, username string, password string, nickname string, email string, mobile string, creator string, status int) (int64, error)
		MemberUpdate(ctx context.Context, id int64, password string, nickname string, email string, mobile string, updater string, status int) error
		MemberDelete(ctx context.Context, id int64, updater string) error
		MemberProfile(ctx context.Context, memberId int64) (*sysentity.SysMember, error)
		FindByUsername(ctx context.Context, username string) (*sysentity.SysMember, error)
	}
)

var (
	localMember IMember
)

func Member() IMember {
	if localMember == nil {
		panic("implement not found for interface IMember, forgot register?")
	}
	return localMember
}

func RegisterMember(i IMember) {
	localMember = i
}
