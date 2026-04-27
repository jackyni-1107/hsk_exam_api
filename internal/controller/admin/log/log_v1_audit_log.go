package log

import (
	"context"

	v1 "exam/api/admin/log/v1"
	syslogsvc "exam/internal/service/syslog"
	"exam/internal/utility"
)

func (c *ControllerV1) AuditLogList(ctx context.Context, req *v1.AuditLogListReq) (res *v1.AuditLogListRes, err error) {
	logs, total, err := syslogsvc.SysLog().AuditLogList(ctx, req.Page, req.Size, req.Username, req.Path, req.Action, req.LogType, req.TraceId, req.StartTime, req.EndTime)
	if err != nil {
		return nil, err
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
		item.CreateTime = utility.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}
	return &v1.AuditLogListRes{List: list, Total: total}, nil
}

func (c *ControllerV1) AuditLogChangeDetails(ctx context.Context, req *v1.AuditLogChangeDetailsReq) (res *v1.AuditLogChangeDetailsRes, err error) {
	details, err := syslogsvc.SysLog().AuditLogChangeDetails(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	list := make([]*v1.AuditChangeDetailItem, 0, len(details))
	for _, d := range details {
		item := &v1.AuditChangeDetailItem{
			Id:          d.Id,
			TableName:   d.TableName,
			RecordId:    d.RecordId,
			FieldName:   d.FieldName,
			BeforeValue: d.BeforeValue,
			AfterValue:  d.AfterValue,
		}
		item.CreateTime = utility.ToRFC3339UTC(d.CreateTime)
		list = append(list, item)
	}
	return &v1.AuditLogChangeDetailsRes{List: list}, nil
}
