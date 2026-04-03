// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// SysConfig is the golang structure of table sys_config for DAO operations like Where/Data.
type SysConfig struct {
	g.Meta      `orm:"table:sys_config, do:true"`
	Id          any         // 主键ID
	ConfigKey   any         // 参数键
	ConfigValue any         // 参数值
	ConfigType  any         // 类型：string/number/boolean/json
	GroupName   any         // 分组
	Remark      any         // 备注
	Creator     any         // 创建者
	CreateTime  *gtime.Time // 创建时间
	Updater     any         // 更新者
	UpdateTime  *gtime.Time // 更新时间
	DeleteFlag  any         // 逻辑删除
}
