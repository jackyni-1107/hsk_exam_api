// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package sys

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SysConfig is the golang structure for table sys_config.
type SysConfig struct {
	Id          int64       `json:"id"           orm:"id"           description:"主键ID"`                          // 主键ID
	ConfigKey   string      `json:"config_key"   orm:"config_key"   description:"参数键"`                           // 参数键
	ConfigValue string      `json:"config_value" orm:"config_value" description:"参数值"`                           // 参数值
	ConfigType  string      `json:"config_type"  orm:"config_type"  description:"类型：string/number/boolean/json"` // 类型：string/number/boolean/json
	GroupName   string      `json:"group_name"   orm:"group_name"   description:"分组"`                            // 分组
	Remark      string      `json:"remark"       orm:"remark"       description:"备注"`                            // 备注
	Creator     string      `json:"creator"      orm:"creator"      description:"创建者"`                           // 创建者
	CreateTime  *gtime.Time `json:"create_time"  orm:"create_time"  description:"创建时间"`                          // 创建时间
	Updater     string      `json:"updater"      orm:"updater"      description:"更新者"`                           // 更新者
	UpdateTime  *gtime.Time `json:"update_time"  orm:"update_time"  description:"更新时间"`                          // 更新时间
	DeleteFlag  int         `json:"delete_flag"  orm:"delete_flag"  description:"逻辑删除"`                          // 逻辑删除
}
