// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttemptCheatEvent is the golang structure of table exam_attempt_cheat_event for DAO operations like Where/Data.
type ExamAttemptCheatEvent struct {
	g.Meta      `orm:"table:exam_attempt_cheat_event, do:true"`
	Id          any         // 主键
	AttemptId   any         // exam_attempt.id
	MemberId    any         // sys_member.id
	EventType   any         // 作弊事件类型，如 switch_screen/screen_record
	EventAt     *gtime.Time // 事件发生时间（服务端写入，与请求到达时刻一致）
	SegmentCode any         // 发生时环节编码 listen/read/write
	Detail      any         // 事件详情
	ClientIp    any         // 客户端IP
	ClientAgent any         // User-Agent
	Creator     any         // 创建者
	CreateTime  *gtime.Time // 创建时间
	DeleteFlag  any         // 逻辑删除
}
