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
	ISysNotification interface {
		ChannelConfigList(ctx context.Context, channel string) ([]sysentity.SysNotificationChannelConfig, error)
		ChannelConfigCreate(ctx context.Context, channel string, provider string, name string, configJson string, creator string) (int64, error)
		ChannelConfigUpdate(ctx context.Context, id int64, name string, configJson string, updater string) error
		ChannelConfigDelete(ctx context.Context, id int64, updater string) error
		ChannelConfigSetActive(ctx context.Context, id int64, updater string) error
		LogList(ctx context.Context, page int, size int, channel string, recipient string) ([]sysentity.SysNotificationLog, int, error)
		Send(ctx context.Context, templateCode string, channel string, recipient string, variables string) (bool, error)
		TemplateList(ctx context.Context, page int, size int, code string, channel string) ([]sysentity.SysNotificationTemplate, int, error)
		TemplateCreate(ctx context.Context, code string, name string, channel string, channelConfigId int64, templateType int, content string, thirdPartyTemplateId string, thirdPartyTemplateParams string, variables string, creator string, status int) (int64, error)
		TemplateUpdate(ctx context.Context, id int64, name string, channel string, channelConfigId int64, templateType int, content string, thirdPartyTemplateId string, thirdPartyTemplateParams string, variables string, updater string, status int) error
		TemplateDelete(ctx context.Context, id int64, updater string) error
	}
)

var (
	localSysNotification ISysNotification
)

func SysNotification() ISysNotification {
	if localSysNotification == nil {
		panic("implement not found for interface ISysNotification, forgot register?")
	}
	return localSysNotification
}

func RegisterSysNotification(i ISysNotification) {
	localSysNotification = i
}
