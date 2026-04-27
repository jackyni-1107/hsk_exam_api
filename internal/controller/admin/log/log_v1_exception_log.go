package log

import (
	"context"

	v1 "exam/api/admin/log/v1"
	syslogsvc "exam/internal/service/syslog"
	"exam/internal/utility"
)

func (c *ControllerV1) ExceptionLogList(ctx context.Context, req *v1.ExceptionLogListReq) (res *v1.ExceptionLogListRes, err error) {
	logs, total, err := syslogsvc.SysLog().ExceptionLogList(ctx, req.Page, req.Size, req.TraceId, req.Path, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.ExceptionLogItem, 0, len(logs))
	for _, e := range logs {
		item := &v1.ExceptionLogItem{
			Id:       e.Id,
			TraceId:  e.TraceId,
			Path:     e.Path,
			Method:   e.Method,
			ErrorMsg: e.ErrorMsg,
			Stack:    e.Stack,
			UserId:   e.UserId,
			Ip:       e.Ip,
		}
		item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}
	return &v1.ExceptionLogListRes{List: list, Total: total}, nil
}
