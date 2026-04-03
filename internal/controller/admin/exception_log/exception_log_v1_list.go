package exception_log

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"exam/api/admin/exception_log/v1"
	"exam/internal/consts"
	daosys "exam/internal/dao/sys"
	entitysys "exam/internal/model/entity/sys"
	"exam/internal/util"
)

func (c *ControllerV1) ExceptionLogList(ctx context.Context, req *v1.ExceptionLogListReq) (res *v1.ExceptionLogListRes, err error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}

	model := daosys.SysExceptionLog.Ctx(ctx)
	if req.TraceId != "" {
		model = model.Where("trace_id", req.TraceId)
	}
	if req.Path != "" {
		model = model.WhereLike("path", "%"+req.Path+"%")
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

	var logs []entitysys.SysExceptionLog
	err = model.Page(req.Page, req.Size).OrderDesc("id").Scan(&logs)
	if err != nil {
		return nil, gerror.WrapCode(consts.CodeInvalidParams, err, "")
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
		item.CreateTime = util.ToRFC3339UTC(e.CreateTime)
		list = append(list, item)
	}

	return &v1.ExceptionLogListRes{List: list, Total: total}, nil
}
