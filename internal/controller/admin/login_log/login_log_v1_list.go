package login_log

import (
	"context"

	v1 "exam/api/admin/login_log/v1"
	syslogsvc "exam/internal/service/syslog"
	"exam/internal/utility"
)

func (c *ControllerV1) LoginLogList(ctx context.Context, req *v1.LoginLogListReq) (res *v1.LoginLogListRes, err error) {
	logs, total, err := syslogsvc.SysLog().LoginLogList(ctx, req.Page, req.Size, req.Username, req.LogType, req.UserType, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.LoginLogItem, 0, len(logs))
	for _, e := range logs {
		item := &v1.LoginLogItem{
			Id:         e.Id,
			LogType:    e.LogType,
			UserId:     e.UserId,
			Username:   e.Username,
			UserType:   e.UserType,
			Ip:         e.Ip,
			UserAgent:  e.UserAgent,
			DeviceInfo: e.DeviceInfo,
			TraceId:    e.TraceId,
			FailReason: e.FailReason,
		}
		item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}
	return &v1.LoginLogListRes{List: list, Total: total}, nil
}
