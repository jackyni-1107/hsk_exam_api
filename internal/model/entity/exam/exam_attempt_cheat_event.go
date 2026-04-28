// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package exam

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ExamAttemptCheatEvent is the golang structure for table exam_attempt_cheat_event.
type ExamAttemptCheatEvent struct {
	Id          int64       `json:"id"           orm:"id"           description:"主键"`                                   // 主键
	AttemptId   int64       `json:"attempt_id"   orm:"attempt_id"   description:"exam_attempt.id"`                      // exam_attempt.id
	MemberId    int64       `json:"member_id"    orm:"member_id"    description:"sys_member.id"`                        // sys_member.id
	EventType   string      `json:"event_type"   orm:"event_type"   description:"作弊事件类型，如 switch_screen/screen_record"` // 作弊事件类型，如 switch_screen/screen_record
	EventAt     *gtime.Time `json:"event_at"     orm:"event_at"     description:"事件发生时间（服务端写入，与请求到达时刻一致）"`              // 事件发生时间（服务端写入，与请求到达时刻一致）
	SegmentCode string      `json:"segment_code" orm:"segment_code" description:"发生时环节编码 listen/read/write"`            // 发生时环节编码 listen/read/write
	Detail      string      `json:"detail"       orm:"detail"       description:"事件详情"`                                 // 事件详情
	ClientIp    string      `json:"client_ip"    orm:"client_ip"    description:"客户端IP"`                                // 客户端IP
	ClientAgent string      `json:"client_agent" orm:"client_agent" description:"User-Agent"`                           // User-Agent
	Creator     string      `json:"creator"      orm:"creator"      description:"创建者"`                                  // 创建者
	CreateTime  *gtime.Time `json:"create_time"  orm:"create_time"  description:"创建时间"`                                 // 创建时间
	DeleteFlag  int         `json:"delete_flag"  orm:"delete_flag"  description:"逻辑删除"`                                 // 逻辑删除
}
