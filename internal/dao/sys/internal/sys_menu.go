// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// SysMenuDao is the data access object for the table sys_menu.
type SysMenuDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  SysMenuColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// SysMenuColumns defines and stores column names for the table sys_menu.
type SysMenuColumns struct {
	Id            string // 菜单ID
	Name          string // 菜单名称
	Permission    string // 权限标识
	Type          string // 菜单类型
	Sort          string // 显示顺序
	ParentId      string // 父菜单ID
	Path          string // 路由地址
	Icon          string // 菜单图标
	Component     string // 组件路径
	ComponentName string // 组件名
	Status        string // 菜单状态
	Visible       string // 是否可见
	KeepAlive     string // 是否缓存
	AlwaysShow    string // 是否总是显示
	Creator       string // 创建者
	CreateTime    string // 创建时间
	Updater       string // 更新者
	UpdateTime    string // 更新时间
	DeleteFlag    string // 逻辑删除标识：0-未删除，1-已删除
}

// sysMenuColumns holds the columns for the table sys_menu.
var sysMenuColumns = SysMenuColumns{
	Id:            "id",
	Name:          "name",
	Permission:    "permission",
	Type:          "type",
	Sort:          "sort",
	ParentId:      "parent_id",
	Path:          "path",
	Icon:          "icon",
	Component:     "component",
	ComponentName: "component_name",
	Status:        "status",
	Visible:       "visible",
	KeepAlive:     "keep_alive",
	AlwaysShow:    "always_show",
	Creator:       "creator",
	CreateTime:    "create_time",
	Updater:       "updater",
	UpdateTime:    "update_time",
	DeleteFlag:    "delete_flag",
}

// NewSysMenuDao creates and returns a new DAO object for table data access.
func NewSysMenuDao(handlers ...gdb.ModelHandler) *SysMenuDao {
	return &SysMenuDao{
		group:    "default",
		table:    "sys_menu",
		columns:  sysMenuColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *SysMenuDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *SysMenuDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *SysMenuDao) Columns() SysMenuColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *SysMenuDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *SysMenuDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *SysMenuDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
