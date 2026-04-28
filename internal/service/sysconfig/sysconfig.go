// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package sysconfig

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysConfig interface {
		ConfigList(ctx context.Context, page int, size int, group string, key string) ([]sysentity.SysConfig, int, error)
		ConfigCreate(ctx context.Context, configKey string, configValue string, configType string, groupName string, remark string, creator string) (int64, error)
		ConfigUpdate(ctx context.Context, id int64, configValue string, remark string, updater string) error
		ConfigDelete(ctx context.Context, id int64, updater string) error
		ConfigGet(ctx context.Context, key string) (string, error)
		ConfigBatchGet(ctx context.Context, keys []string) (map[string]string, error)
	}
)

var (
	localSysConfig ISysConfig
)

func SysConfig() ISysConfig {
	if localSysConfig == nil {
		panic("implement not found for interface ISysConfig, forgot register?")
	}
	return localSysConfig
}

func RegisterSysConfig(i ISysConfig) {
	localSysConfig = i
}
