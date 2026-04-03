// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysAuditChangeDetail is the golang structure of table sys_audit_change_detail for DAO operations like Where/Data.
type SysAuditChangeDetail struct {
	g.Meta         `orm:"table:sys_audit_change_detail, do:true"`
	Id             any         // 主键ID
	OperationLogId any         // 关联system_operation_audit_log.id
	TableName      any         // 表名
	RecordId       any         // 记录ID
	FieldName      any         // 字段名
	BeforeValue    any         // 变更前值
	AfterValue     any         // 变更后值
	CreateTime     *gtime.Time // 创建时间
}
