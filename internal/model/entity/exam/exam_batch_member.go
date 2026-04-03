// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamBatchMember is the golang structure for table exam_batch_member.
type ExamBatchMember struct {
	BatchId    int64       `json:"batch_id"    orm:"batch_id"    description:"exam_batch.id"` // exam_batch.id
	MemberId   int64       `json:"member_id"   orm:"member_id"   description:"sys_member.id"` // sys_member.id
	Creator    string      `json:"creator"     orm:"creator"     description:"导入操作者"`         // 导入操作者
	CreateTime *gtime.Time `json:"create_time" orm:"create_time" description:"导入时间"`          // 导入时间
}
