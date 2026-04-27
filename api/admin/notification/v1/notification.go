package v1

import "github.com/gogf/gf/v2/frame/g"

type LogListReq struct {
	g.Meta    `path:"/notification/log/list" method:"get" tags:"通知" summary:"发送记录" permission:"notification:log_list"`
	Page      int    `json:"page" dc:"页码"`
	Size      int    `json:"size" dc:"每页条数"`
	Channel   string `json:"channel" dc:"渠道"`
	Recipient string `json:"recipient" dc:"接收人"`
}

type LogListRes struct {
	List  []*LogItem `json:"list" dc:"列表"`
	Total int        `json:"total" dc:"总数"`
}

type LogItem struct {
	Id           int64  `json:"id" dc:"记录ID"`
	TemplateCode string `json:"template_code" dc:"模板编码"`
	Channel      string `json:"channel" dc:"渠道"`
	Recipient    string `json:"recipient" dc:"接收人"`
	Status       int    `json:"status" dc:"发送状态"`
	ErrorMsg     string `json:"error_msg" dc:"错误信息"`
	CreateTime   string `json:"create_time" dc:"创建时间"`
}

type ChannelConfigListReq struct {
	g.Meta  `path:"/notification/channel/list" method:"get" tags:"通知" summary:"渠道配置列表" permission:"notification:channel_config_list"`
	Channel string `json:"channel" dc:"渠道"`
}

type ChannelConfigListRes struct {
	List []*ChannelConfigItem `json:"list" dc:"渠道配置列表"`
}

type ChannelConfigItem struct {
	Id         int64  `json:"id" dc:"配置ID"`
	Channel    string `json:"channel" dc:"渠道"`
	Provider   string `json:"provider" dc:"服务商"`
	Name       string `json:"name" dc:"配置名称"`
	IsActive   int    `json:"is_active" dc:"是否当前使用：0否 1是"`
	ConfigJson string `json:"config_json" dc:"配置JSON"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type ChannelConfigCreateReq struct {
	g.Meta     `path:"/notification/channel" method:"post" tags:"通知" summary:"新增渠道配置" permission:"notification:channel_config_create"`
	Channel    string `json:"channel" v:"required#err.invalid_params" dc:"渠道"`
	Provider   string `json:"provider" v:"required#err.invalid_params" dc:"服务商"`
	Name       string `json:"name" v:"required#err.invalid_params" dc:"配置名称"`
	ConfigJson string `json:"config_json" v:"required#err.invalid_params" dc:"配置JSON"`
}

type ChannelConfigCreateRes struct {
	Id int64 `json:"id" dc:"配置ID"`
}

type ChannelConfigUpdateReq struct {
	g.Meta     `path:"/notification/channel/{id}" method:"put" tags:"通知" summary:"更新渠道配置" permission:"notification:channel_config_update"`
	Id         int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
	Name       string `json:"name" dc:"配置名称"`
	ConfigJson string `json:"config_json" dc:"配置JSON"`
}

type ChannelConfigUpdateRes struct{}

type ChannelConfigDeleteReq struct {
	g.Meta `path:"/notification/channel/{id}" method:"delete" tags:"通知" summary:"删除渠道配置" permission:"notification:channel_config_delete"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
}

type ChannelConfigDeleteRes struct{}

type ChannelConfigSetActiveReq struct {
	g.Meta `path:"/notification/channel/{id}/set-active" method:"post" tags:"通知" summary:"设为当前渠道" permission:"notification:channel_config_set_active"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"配置ID"`
}

type ChannelConfigSetActiveRes struct{}

type TemplateListReq struct {
	g.Meta  `path:"/notification/template/list" method:"get" tags:"通知" summary:"模板列表" permission:"notification:template_list"`
	Page    int    `json:"page" dc:"页码"`
	Size    int    `json:"size" dc:"每页条数"`
	Code    string `json:"code" dc:"模板编码"`
	Channel string `json:"channel" dc:"渠道"`
}

type TemplateListRes struct {
	List  []*TemplateItem `json:"list" dc:"列表"`
	Total int             `json:"total" dc:"总数"`
}

type TemplateItem struct {
	Id         int64  `json:"id" dc:"模板ID"`
	Code       string `json:"code" dc:"模板编码"`
	Name       string `json:"name" dc:"模板名称"`
	Channel    string `json:"channel" dc:"渠道"`
	Content    string `json:"content" dc:"模板内容"`
	Variables  string `json:"variables" dc:"变量(JSON)"`
	Status     int    `json:"status" dc:"状态：0启用 1停用"`
	CreateTime string `json:"create_time" dc:"创建时间"`
}

type TemplateCreateReq struct {
	g.Meta    `path:"/notification/template" method:"post" tags:"通知" summary:"新增模板" permission:"notification:template_create"`
	Code      string `json:"code" v:"required#err.invalid_params" dc:"模板编码"`
	Name      string `json:"name" v:"required#err.invalid_params" dc:"模板名称"`
	Channel   string `json:"channel" v:"required#err.invalid_params" dc:"渠道"`
	Content   string `json:"content" v:"required#err.invalid_params" dc:"模板内容"`
	Variables string `json:"variables" dc:"变量(JSON)"`
	Status    int    `json:"status" dc:"状态：0启用 1停用"`
}

type TemplateCreateRes struct {
	Id int64 `json:"id" dc:"模板ID"`
}

type TemplateUpdateReq struct {
	g.Meta    `path:"/notification/template/{id}" method:"put" tags:"通知" summary:"更新模板" permission:"notification:template_update"`
	Id        int64  `json:"id" in:"path" v:"required#err.invalid_params" dc:"模板ID"`
	Name      string `json:"name" dc:"模板名称"`
	Content   string `json:"content" dc:"模板内容"`
	Variables string `json:"variables" dc:"变量(JSON)"`
	Status    int    `json:"status" dc:"状态：0启用 1停用"`
}

type TemplateUpdateRes struct{}

type TemplateDeleteReq struct {
	g.Meta `path:"/notification/template/{id}" method:"delete" tags:"通知" summary:"删除模板" permission:"notification:template_delete"`
	Id     int64 `json:"id" in:"path" v:"required#err.invalid_params" dc:"模板ID"`
}

type TemplateDeleteRes struct{}

type NotificationSendReq struct {
	g.Meta       `path:"/notification/send" method:"post" tags:"通知" summary:"发送通知" permission:"notification:send"`
	TemplateCode string `json:"template_code" v:"required#err.invalid_params" dc:"模板编码"`
	Channel      string `json:"channel" v:"required#err.invalid_params" dc:"渠道"`
	Recipient    string `json:"recipient" v:"required#err.invalid_params" dc:"接收人"`
	Variables    string `json:"variables" dc:"变量(JSON)"`
}

type NotificationSendRes struct {
	Ok bool `json:"ok" dc:"是否成功"`
}
