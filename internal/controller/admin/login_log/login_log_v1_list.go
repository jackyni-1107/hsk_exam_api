package login_log

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	v1 "exam/api/admin/login_log/v1"
	"exam/internal/consts"
	sysdao "exam/internal/dao/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/util"
)

func (c *ControllerV1) LoginLogList(ctx context.Context, req *v1.LoginLogListReq) (res *v1.LoginLogListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	model := sysdao.SysLoginLog.Ctx(ctx)
	if req.Username != "" {
		model = model.WhereLike("username", "%"+req.Username+"%")
	}
	if req.LogType != "" {
		model = model.Where("log_type", req.LogType)
	}
	if req.UserType >= 1 {
		model = model.Where("user_type", req.UserType)
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

	var logs []sysentity.SysLoginLog
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&logs)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
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
		item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}

	return &v1.LoginLogListRes{List: list, Total: total}, nil
}
