package config

import (
	"context"

	v1 "exam/api/client/config/v1"
	sysconfigsvc "exam/internal/service/sysconfig"
)

func (c *ControllerV1) MultiGet(ctx context.Context, req *v1.MultiGetReq) (res *v1.MultiGetRes, err error) {
	kvMap, err := sysconfigsvc.SysConfig().ConfigBatchGet(ctx, req.Keys)
	if err != nil {
		return nil, err
	}
	items := make([]*v1.ConfigKV, 0, len(req.Keys))
	for _, key := range req.Keys {
		items = append(items, &v1.ConfigKV{
			Key:   key,
			Value: kvMap[key],
		})
	}
	return &v1.MultiGetRes{Items: items}, nil
}
