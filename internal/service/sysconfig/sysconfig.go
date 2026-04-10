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
		ConfigList(ctx context.Context, page, size int, group, key string) ([]sysentity.SysConfig, int, error)
		ConfigCreate(ctx context.Context, configKey, configValue, configType, groupName, remark, creator string) (int64, error)
		ConfigUpdate(ctx context.Context, id int64, configValue, remark, updater string) error
		ConfigDelete(ctx context.Context, id int64, updater string) error
		ConfigGet(ctx context.Context, key string) (string, error)
		DictTypeList(ctx context.Context, page, size int, dictType string) ([]sysentity.SysDictType, int, error)
		DictTypeCreate(ctx context.Context, dictName, dictType, remark, creator string, status int) (int64, error)
		DictTypeUpdate(ctx context.Context, id int64, dictName, remark, updater string, status int) error
		DictTypeDelete(ctx context.Context, id int64, updater string) error
		DictDataList(ctx context.Context, page, size int, dictType string) ([]sysentity.SysDictData, int, error)
		DictDataCreate(ctx context.Context, dictType, dictLabel, dictValue, remark, creator string, sort, status int) (int64, error)
		DictDataUpdate(ctx context.Context, id int64, dictLabel, dictValue, remark, updater string, sort, status int) error
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
