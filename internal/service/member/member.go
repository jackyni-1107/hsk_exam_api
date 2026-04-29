// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package member

import (
	"context"
	"io"

	"exam/internal/model/bo"
	sysentity "exam/internal/model/entity/sys"
)

type (
	IMember interface {
		MemberList(ctx context.Context, page int, size int, username string, status int) ([]sysentity.SysMember, int, error)
		MemberCreate(ctx context.Context, username string, password string, nickname string, email string, mobile string, creator string, status int) (int64, error)
		MemberUpdate(ctx context.Context, id int64, password string, nickname string, email string, mobile string, updater string, status int) error
		MemberDelete(ctx context.Context, id int64, updater string) error
		// MemberImport 从 CSV 流批量创建客户（UTF-8，首行为表头）。country/year/seqDigits 由调用方传入，用于生成用户名。
		MemberImport(ctx context.Context, r io.Reader, creator string, country, year string, seqDigits int) (*bo.MemberImportResult, error)
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
