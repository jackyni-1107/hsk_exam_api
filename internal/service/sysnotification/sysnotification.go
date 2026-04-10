// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package sysnotification

import (
	"context"
	sysentity "exam/internal/model/entity/sys"
)

type (
	ISysnotification interface {
		TemplateList(ctx context.Context, page, size int, code, channel string) ([]sysentity.SysNotificationTemplate, int, error)
		TemplateCreate(ctx context.Context, code, name, channel, content, variables, creator string, status int) (int64, error)
		TemplateUpdate(ctx context.Context, id int64, name, content, variables, updater string, status int) error
		TemplateDelete(ctx context.Context, id int64, updater string) error
		ChannelConfigList(ctx context.Context, channel string) ([]sysentity.SysNotificationChannelConfig, error)
		ChannelConfigCreate(ctx context.Context, channel, provider, name, configJson, creator string) (int64, error)
		ChannelConfigUpdate(ctx context.Context, id int64, name, configJson, updater string) error
		ChannelConfigDelete(ctx context.Context, id int64, updater string) error
		ChannelConfigSetActive(ctx context.Context, id int64, updater string) error
		Send(ctx context.Context, templateCode, channel, recipient, variables string) (bool, error)
		LogList(ctx context.Context, page, size int, channel, recipient string) ([]sysentity.SysNotificationLog, int, error)
	}
)

var (
	localSysnotification ISysnotification
)

func Sysnotification() ISysnotification {
	if localSysnotification == nil {
		panic("implement not found for interface ISysnotification, forgot register?")
	}
	return localSysnotification
}

func RegisterSysnotification(i ISysnotification) {
	localSysnotification = i
}
