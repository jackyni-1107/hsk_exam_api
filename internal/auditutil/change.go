package auditutil

import (
	"context"

	"exam/internal/middleware"
	auditsvc "exam/internal/service/audit"

	"github.com/gogf/gf/v2/util/gconv"
)

// RecordEntityDiff 在存在 operation_log_id 时，将实体前后差异写入 sys_audit_change_detail。
// before/after 可为 nil；未走 Audit 中间件的请求会自动跳过。
func RecordEntityDiff(ctx context.Context, tableName string, recordId int64, before, after interface{}) {
	opID := middleware.GetOperationLogId(ctx)
	if opID <= 0 {
		return
	}
	auditsvc.Audit().RecordChange(ctx, tableName, recordId, opID, entityToMap(before), entityToMap(after))
}

// RecordMapDiff 使用已构造的字段 map 记录差异（如关联表 role_ids 摘要）。
func RecordMapDiff(ctx context.Context, tableName string, recordId int64, before, after map[string]interface{}) {
	opID := middleware.GetOperationLogId(ctx)
	if opID <= 0 {
		return
	}
	if before == nil {
		before = map[string]interface{}{}
	}
	if after == nil {
		after = map[string]interface{}{}
	}
	auditsvc.Audit().RecordChange(ctx, tableName, recordId, opID, before, after)
}

func entityToMap(v interface{}) map[string]interface{} {
	if v == nil {
		return nil
	}
	m := gconv.Map(v)
	if m == nil {
		return map[string]interface{}{}
	}
	return m
}
