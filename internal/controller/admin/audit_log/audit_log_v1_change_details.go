package audit_log

import (
	"context"

	v1 "exam/api/admin/audit_log/v1"
	syslogsvc "exam/internal/service/syslog"
	"exam/internal/utility"
)

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
