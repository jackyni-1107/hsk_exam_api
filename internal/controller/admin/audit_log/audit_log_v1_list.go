package audit_log

import (
	"context"

	"exam/api/admin/audit_log/v1"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/model/entity"
	"exam/internal/util"
	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) AuditLogList(ctx context.Context, req *v1.AuditLogListReq) (res *v1.AuditLogListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	model := dao.SystemOperationAuditLog.Ctx(ctx)
	if req.Username != "" {
		model = model.WhereLike("username", "%"+req.Username+"%")
	}
	if req.Path != "" {
		model = model.WhereLike("path", "%"+req.Path+"%")
	}
	if req.Action != "" {
		model = model.Where("action", req.Action)
	}
	if req.LogType != "" {
		model = model.Where("log_type", req.LogType)
	}
	if req.TraceId != "" {
		model = model.Where("trace_id", req.TraceId)
	}
	if req.StartTime != "" {
		model = model.WhereGTE("create_time", req.StartTime)
	}
	if req.EndTime != "" {
		model = model.WhereLTE("create_time", req.EndTime)
	}

	total, err := model.Count()
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}

	var logs []entity.SystemOperationAuditLog
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&logs)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}

	list := make([]*v1.AuditLogItem, 0, len(logs))
	for _, e := range logs {
		item := &v1.AuditLogItem{
			Id:           int64(e.Id),
			UserId:       e.UserId,
			Username:     e.Username,
			UserType:     e.UserType,
			Module:       e.Module,
			Action:       e.Action,
			LogType:      e.LogType,
			Method:       e.Method,
			Path:         e.Path,
			RequestData:  e.RequestData,
			ResponseData: e.ResponseData,
			Ip:           e.Ip,
			UserAgent:    e.UserAgent,
			TraceId:      e.TraceId,
			DeviceInfo:   e.DeviceInfo,
			DurationMs:   e.DurationMs,
		}
		item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}

	return &v1.AuditLogListRes{List: list, Total: total}, nil
}
