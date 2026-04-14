// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package batch

import (
	"context"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	IBatch interface {
		// ExamBatchList 分页查询考试批次列表
		ExamBatchList(ctx context.Context, page int, size int, key string) (list []examentity.ExamBatch, total int, err error)
		// ExamBatchDetail 批次详情（含 Mock 卷 id 列表与学员数）。
		ExamBatchDetail(ctx context.Context, id int64) (*bo.ExamBatchAdminItem, error)
		// ExamBatchCreate 创建考试批次
		ExamBatchCreate(ctx context.Context, name string, startAt string, endAt string, creator string) (int64, error)
		// ExamBatchUpdate 更新考试批次
		ExamBatchUpdate(ctx context.Context, id int64, name string, startAt string, endAt string, updater string) error
		// ExamBatchDelete 删除考试批次
		ExamBatchDelete(ctx context.Context, id int64, updater string) error
		// ExamBatchMembersAdd 批量向指定批次和 Mock 卷添加学员
		ExamBatchMembersAdd(ctx context.Context, batchID int64, mockPaperID int64, memberIDs []int64, creator string) error
		// ExamBatchMembersRemove 从批次中移除学员
		ExamBatchMembersRemove(ctx context.Context, batchID int64, mockPaperID int64, memberIDs []int64) (int, error)
		// ExamBatchMemberList 查询批次内的成员列表（关联系统用户表）
		ExamBatchMemberList(ctx context.Context, page int, size int, batchID int64, mockPaperID int64) (list []sysentity.SysUser, total int, err error)
		// MyExamBatches 学员查询自己的批次
		MyExamBatches(ctx context.Context, memberID int64) (list []bo.MyExamBatchItem, err error)
		// GetExamWindowStatus 判定考试窗口状态
		GetExamWindowStatus(now *gtime.Time, start *gtime.Time, end *gtime.Time) string
		// GetBatchByID 内部公用查询
		GetBatchByID(ctx context.Context, id int64) (*examentity.ExamBatch, error)
	}
)

var (
	localBatch IBatch
)

func Batch() IBatch {
	if localBatch == nil {
		panic("implement not found for interface IBatch, forgot register?")
	}
	return localBatch
}

func RegisterBatch(i IBatch) {
	localBatch = i
}
