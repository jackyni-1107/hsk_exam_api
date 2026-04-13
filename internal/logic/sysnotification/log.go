package sysnotification

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysNotification) LogList(ctx context.Context, page, size int, channel, recipient string) ([]sysentity.SysNotificationLog, int, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	m := dao.SysNotificationLog.Ctx(ctx)
	if channel != "" {
		m = m.Where("channel", channel)
	}
	if recipient != "" {
		m = m.WhereLike("recipient", "%"+recipient+"%")
	}
	total, err := m.Count()
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	var list []sysentity.SysNotificationLog
	err = m.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, gerror.WrapCode(consts.CodeInvalidParams, err, "")
	}
	return list, total, nil
}
