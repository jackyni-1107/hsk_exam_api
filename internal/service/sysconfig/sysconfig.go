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
	ISysconfig interface {
		ConfigList(ctx context.Context, page int, size int, group string, key string) ([]sysentity.SysConfig, int, error)
		ConfigCreate(ctx context.Context, configKey string, configValue string, configType string, groupName string, remark string, creator string) (int64, error)
		ConfigUpdate(ctx context.Context, id int64, configValue string, remark string, updater string) error
		ConfigDelete(ctx context.Context, id int64, updater string) error
		ConfigGet(ctx context.Context, key string) (string, error)
		DictTypeList(ctx context.Context, page int, size int, dictType string) ([]sysentity.SysDictType, int, error)
		DictTypeCreate(ctx context.Context, dictName string, dictType string, remark string, creator string, status int) (int64, error)
		DictTypeUpdate(ctx context.Context, id int64, dictName string, remark string, updater string, status int) error
		DictTypeDelete(ctx context.Context, id int64, updater string) error
		DictDataList(ctx context.Context, page int, size int, dictType string) ([]sysentity.SysDictData, int, error)
		DictDataCreate(ctx context.Context, dictType string, dictLabel string, dictValue string, remark string, creator string, sort int, status int) (int64, error)
		DictDataUpdate(ctx context.Context, id int64, dictLabel string, dictValue string, remark string, updater string, sort int, status int) error
		DictDataDelete(ctx context.Context, id int64, updater string) error
	}
)

var (
	localSysconfig ISysconfig
)

func Sysconfig() ISysconfig {
	if localSysconfig == nil {
		panic("implement not found for interface ISysconfig, forgot register?")
	}
	return localSysconfig
}

func RegisterSysconfig(i ISysconfig) {
	localSysconfig = i
}
