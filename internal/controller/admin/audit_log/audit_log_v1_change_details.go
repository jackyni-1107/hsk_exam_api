package audit_log

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/admin/audit_log/v1"
	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"
)

func (c *ControllerV1) AuditLogChangeDetails(ctx context.Context, req *v1.AuditLogChangeDetailsReq) (res *v1.AuditLogChangeDetailsRes, err error) {
	var details []sysentity.SysAuditChangeDetail
	err = sysdao.SysAuditChangeDetail.Ctx(ctx).Where("operation_log_id", req.Id).OrderAsc("id").Scan(&details)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
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
		item.CreateTime = util.ToRFC3339UTC(d.CreateTime)
		list = append(list, item)
	}
	return &v1.AuditLogChangeDetailsRes{List: list}, nil
}
