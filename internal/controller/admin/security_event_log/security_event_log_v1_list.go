package security_event_log

import (
	"context"

	v1 "exam/api/admin/security_event_log/v1"
	syslogsvc "exam/internal/service/syslog"
	"exam/internal/utility"
)

func (c *ControllerV1) SecurityEventLogList(ctx context.Context, req *v1.SecurityEventLogListReq) (res *v1.SecurityEventLogListRes, err error) {
	logs, total, err := syslogsvc.SysLog().SecurityEventLogList(ctx, req.Page, req.Size, req.EventType, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
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
		item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}
	return &v1.SecurityEventLogListRes{List: list, Total: total}, nil
}
