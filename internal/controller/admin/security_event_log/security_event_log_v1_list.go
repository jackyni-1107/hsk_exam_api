package security_event_log

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/api/admin/security_event_log/v1"
	"exam/internal/consts"
	daosys "exam/internal/dao/sys"
	entitysys "exam/internal/model/entity/sys"
	"exam/internal/util"
)

func (c *ControllerV1) SecurityEventLogList(ctx context.Context, req *v1.SecurityEventLogListReq) (res *v1.SecurityEventLogListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	model := daosys.SysSecurityEventLog.Ctx(ctx)
	if req.EventType != "" {
		model = model.Where("event_type", req.EventType)
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

	var logs []entitysys.SysSecurityEventLog
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&logs)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}

	list := make([]*v1.SecurityEventLogItem, 0, len(logs))
	for _, e := range logs {
		item := &v1.SecurityEventLogItem{
			Id:        e.Id,
			EventType: e.EventType,
			UserId:    e.UserId,
			Ip:        e.Ip,
			UserAgent: e.UserAgent,
			Detail:    e.Detail,
			TraceId:   e.TraceId,
		}
		item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}

	return &v1.SecurityEventLogListRes{List: list, Total: total}, nil
}
