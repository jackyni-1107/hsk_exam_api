// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package attempt

import (
	"context"
	"exam/internal/model/bo"
	examentity "exam/internal/model/entity/exam"

	"github.com/gogf/gf/v2/os/gtime"
)

type (
	IAttempt interface {
		// AttemptAdminList 分页查询答题会话（联表学员、试卷）。
		AttemptAdminList(ctx context.Context, page int, size int, level string, examinationPaperId int64, examBatchId int64, status int, username string) ([]bo.AttemptAdminListRow, int, error)
		// AttemptAdminDetail 按 id 加载会话、学员、试卷及答题明细（含客观题是否选对）。
		AttemptAdminDetail(ctx context.Context, attemptID int64) (*bo.AttemptAdminDetailView, error)
		// AttemptAdminSaveSubjectiveScores 写入主观题人工分并汇总 subjective_score、total_score（允许部分题目已评）。
		AttemptAdminSaveSubjectiveScores(ctx context.Context, attemptID int64, items []bo.SubjectiveScoreItem) (subjectiveSum float64, totalScore float64, err error)
		// CreateAttemptForBatch 按批次创建会话（未开始）；examPaperID 为多卷批次必选（见 exam_batch_member.exam_paper_id）；单卷批次可传 0。
		CreateAttemptForBatch(ctx context.Context, userID int64, batchID int64, examPaperID int64) (int64, error)
		// StartAttempt 开考：进入进行中并写入截止时间。
		StartAttempt(ctx context.Context, userID int64, attemptID int64, clientDurationSeconds int) error
		// GetAttempt 查询会话；若已超时仍进行中则自动交卷并计分。
		GetAttempt(ctx context.Context, userID int64, attemptID int64) (*bo.AttemptView, error)
		// GetAttemptAnswers 返回当前用户该会话下的答题明细：先读库再合并 Redis 草稿（与保存路径一致），仅包含非空答案。
		GetAttemptAnswers(ctx context.Context, userID int64, attemptID int64) ([]bo.AttemptAnswerClientItem, error)
		// SaveAnswers 保存答案 redis -> db
		SaveAnswers(ctx context.Context, userID int64, attemptID int64, segmentCode string, items []bo.SaveAnswerItem) error
		// SubmitAttempt 主动交卷：仅标记为已交卷（待算分）。客观分与 exam_result 由 sys_task（ExamScoreFinalizeHandler）统一算分写入。
		SubmitAttempt(ctx context.Context, userID int64, attemptID int64) error
		// MarkSubmittedIfOverdue 供定时任务：超时未操作会话标记为已交卷（待算分，不校验用户）。算分由 ExamScoreFinalizeHandler 执行。
		MarkSubmittedIfOverdue(ctx context.Context, attemptID int64) error
		// MarkSubmittedByBatchExpired 供定时任务：批次过期后进行中会话标记为已交卷（待算分，不校验用户）。
		MarkSubmittedByBatchExpired(ctx context.Context, attemptID int64) error
		// FinalizeAttempt 对已交卷（待算分）会话计算客观分并置为已结束，写入 exam_result。仅应由 ExamScoreFinalizeHandler（sys_task）调用。
		FinalizeAttempt(ctx context.Context, attemptID int64) error
		// GetAttemptByID 获取答题会话详情
		GetAttemptByID(ctx context.Context, id int64) (*examentity.ExamAttempt, error)
		// IsWindowOpen 判断考试窗口是否开启（含起止边界时刻）。
		IsWindowOpen(now *gtime.Time, start *gtime.Time, end *gtime.Time) bool
	}
)

var (
	localAttempt IAttempt
)

func Attempt() IAttempt {
	if localAttempt == nil {
		panic("implement not found for interface IAttempt, forgot register?")
	}
	return localAttempt
}

func RegisterAttempt(i IAttempt) {
	localAttempt = i
}
