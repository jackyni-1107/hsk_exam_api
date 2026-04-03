// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAuditChangeDetail is the golang structure for table sys_audit_change_detail.
type SysAuditChangeDetail struct {
	Id             int64       `json:"id"               orm:"id"               description:"主键ID"`                            // 主键ID
	OperationLogId int64       `json:"operation_log_id" orm:"operation_log_id" description:"关联system_operation_audit_log.id"` // 关联system_operation_audit_log.id
	TableName      string      `json:"table_name"       orm:"table_name"       description:"表名"`                              // 表名
	RecordId       int64       `json:"record_id"        orm:"record_id"        description:"记录ID"`                            // 记录ID
	FieldName      string      `json:"field_name"       orm:"field_name"       description:"字段名"`                             // 字段名
	BeforeValue    string      `json:"before_value"     orm:"before_value"     description:"变更前值"`                            // 变更前值
	AfterValue     string      `json:"after_value"      orm:"after_value"      description:"变更后值"`                            // 变更后值
	CreateTime     *gtime.Time `json:"create_time"      orm:"create_time"      description:"创建时间"`                            // 创建时间
}
