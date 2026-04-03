package entity

import (
	entityexam "exam/internal/model/entity/exam"
	entitysys "exam/internal/model/entity/sys"
)

// 与历史代码中的 System* 命名兼容（底层为 sys_* 表）。
type (
	SystemUser              = entitysys.SysUser
	SystemMember            = entitysys.SysMember
	SystemUserRole          = entitysys.SysUserRole
	SystemRole              = entitysys.SysRole
	SystemRoleMenu          = entitysys.SysRoleMenu
	SystemMenu              = entitysys.SysMenu
	SystemConfig            = entitysys.SysConfig
	SystemDictType          = entitysys.SysDictType
	SystemDictData          = entitysys.SysDictData
	SystemOperationAuditLog = entitysys.SysOperationAuditLog

	SysTask                      = entitysys.SysTask
	SysTaskLog                   = entitysys.SysTaskLog
	SysFileStorage               = entitysys.SysFileStorage
	SysFileStorageConfig         = entitysys.SysFileStorageConfig
	SysNotificationLog           = entitysys.SysNotificationLog
	SysNotificationChannelConfig = entitysys.SysNotificationChannelConfig
	SysNotificationTemplate      = entitysys.SysNotificationTemplate
)

// 考试域实体别名
type (
	ExamPaper         = entityexam.ExamPaper
	ExamSection       = entityexam.ExamSection
	ExamQuestionBlock = entityexam.ExamQuestionBlock
	ExamQuestion      = entityexam.ExamQuestion
	ExamOption        = entityexam.ExamOption
	ExamAttempt       = entityexam.ExamAttempt
	ExamAttemptAnswer = entityexam.ExamAttemptAnswer
	ExamResult        = entityexam.ExamResult
)

// ClientUser 客户端学员（数据在 sys_member 表，与 SystemMember 同义）
type ClientUser = entitysys.SysMember
