package v1

import "github.com/gogf/gf/v2/frame/g"

type LogListReq struct {
	g.Meta    `path:"/notification/log/list" method:"get" tags:"通知" summary:"发送记录"`
	Page      int    `json:"page"`
	Size      int    `json:"size"`
	Channel   string `json:"channel"`
	Recipient string `json:"recipient"`
}

type LogListRes struct {
	List  []*LogItem `json:"list"`
	Total int        `json:"total"`
}

type LogItem struct {
	Id           int64  `json:"id"`
	TemplateCode string `json:"template_code"`
	Channel      string `json:"channel"`
	Recipient    string `json:"recipient"`
	Status       int    `json:"status"`
	ErrorMsg     string `json:"error_msg"`
	CreateTime   string `json:"create_time"`
}

type ChannelConfigListReq struct {
	g.Meta  `path:"/notification/channel/list" method:"get" tags:"通知" summary:"渠道配置列表"`
	Channel string `json:"channel"`
}

type ChannelConfigListRes struct {
	List []*ChannelConfigItem `json:"list"`
}

type ChannelConfigItem struct {
	Id         int64  `json:"id"`
	Channel    string `json:"channel"`
	Provider   string `json:"provider"`
	Name       string `json:"name"`
	IsActive   int    `json:"is_active"`
	ConfigJson string `json:"config_json"`
	CreateTime string `json:"create_time"`
}

type ChannelConfigCreateReq struct {
	g.Meta     `path:"/notification/channel" method:"post" tags:"通知" summary:"新增渠道配置"`
	Channel    string `json:"channel" v:"required#err.invalid_params"`
	Provider   string `json:"provider" v:"required#err.invalid_params"`
	Name       string `json:"name" v:"required#err.invalid_params"`
	ConfigJson string `json:"config_json" v:"required#err.invalid_params"`
}

type ChannelConfigCreateRes struct {
	Id int64 `json:"id"`
}

type ChannelConfigUpdateReq struct {
	g.Meta     `path:"/notification/channel/{id}" method:"put" tags:"通知" summary:"更新渠道配置"`
	Id         int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	Name       string `json:"name"`
	ConfigJson string `json:"config_json"`
}

type ChannelConfigUpdateRes struct{}

type ChannelConfigDeleteReq struct {
	g.Meta `path:"/notification/channel/{id}" method:"delete" tags:"通知" summary:"删除渠道配置"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type ChannelConfigDeleteRes struct{}

type ChannelConfigSetActiveReq struct {
	g.Meta `path:"/notification/channel/{id}/set-active" method:"post" tags:"通知" summary:"设为当前渠道"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type ChannelConfigSetActiveRes struct{}

type TemplateListReq struct {
	g.Meta  `path:"/notification/template/list" method:"get" tags:"通知" summary:"模板列表"`
	Page    int    `json:"page"`
	Size    int    `json:"size"`
	Code    string `json:"code"`
	Channel string `json:"channel"`
}

type TemplateListRes struct {
	List  []*TemplateItem `json:"list"`
	Total int             `json:"total"`
}

type TemplateItem struct {
	Id         int64  `json:"id"`
	Code       string `json:"code"`
	Name       string `json:"name"`
	Channel    string `json:"channel"`
	Content    string `json:"content"`
	Variables  string `json:"variables"`
	Status     int    `json:"status"`
	CreateTime string `json:"create_time"`
}

type TemplateCreateReq struct {
	g.Meta    `path:"/notification/template" method:"post" tags:"通知" summary:"新增模板"`
	Code      string `json:"code" v:"required#err.invalid_params"`
	Name      string `json:"name" v:"required#err.invalid_params"`
	Channel   string `json:"channel" v:"required#err.invalid_params"`
	Content   string `json:"content" v:"required#err.invalid_params"`
	Variables string `json:"variables"`
	Status    int    `json:"status"`
}

type TemplateCreateRes struct {
	Id int64 `json:"id"`
}

type TemplateUpdateReq struct {
	g.Meta    `path:"/notification/template/{id}" method:"put" tags:"通知" summary:"更新模板"`
	Id        int64  `json:"id" in:"path" v:"required#err.invalid_params"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	Variables string `json:"variables"`
	Status    int    `json:"status"`
}

type TemplateUpdateRes struct{}

type TemplateDeleteReq struct {
	g.Meta `path:"/notification/template/{id}" method:"delete" tags:"通知" summary:"删除模板"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params"`
}

type TemplateDeleteRes struct{}

type NotificationSendReq struct {
	g.Meta       `path:"/notification/send" method:"post" tags:"通知" summary:"发送通知"`
	TemplateCode string `json:"template_code" v:"required#err.invalid_params"`
	Channel      string `json:"channel" v:"required#err.invalid_params"`
	Recipient    string `json:"recipient" v:"required#err.invalid_params"`
	Variables    string `json:"variables"`
}

type NotificationSendRes struct {
	Ok bool `json:"ok"`
}
