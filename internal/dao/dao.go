package dao

import (
	examdao "exam/internal/dao/exam"
	mockdao "exam/internal/dao/mock"
	sysdao "exam/internal/dao/sys"
)

// 管理端 RBAC 等使用的 System* 命名（底层为 sys_* DAO）。
var (
	SystemUser     = sysdao.SysUser
	SystemUserRole = sysdao.SysUserRole
	SystemRoleMenu = sysdao.SysRoleMenu
	SystemMenu     = sysdao.SysMenu
	SystemRole     = sysdao.SysRole

	SystemDictType = sysdao.SysDictType
	SystemDictData = sysdao.SysDictData
	SystemConfig   = sysdao.SysConfig

	SystemOperationAuditLog = sysdao.SysOperationAuditLog
)

// Sys* 表（任务、文件、通知等）
var (
	SysTask                      = sysdao.SysTask
	SysTaskLog                   = sysdao.SysTaskLog
	SysFileStorage               = sysdao.SysFileStorage
	SysFileStorageConfig         = sysdao.SysFileStorageConfig
	SysNotificationChannelConfig = sysdao.SysNotificationChannelConfig
	SysNotificationTemplate      = sysdao.SysNotificationTemplate
	SysNotificationLog           = sysdao.SysNotificationLog
)

// 考试与模拟卷
var (
	ExamPaper                = examdao.ExamPaper
	ExamSection              = examdao.ExamSection
	ExamQuestionBlock        = examdao.ExamQuestionBlock
	ExamQuestion             = examdao.ExamQuestion
	ExamOption               = examdao.ExamOption
	ExamAttempt              = examdao.ExamAttempt
	ExamAttemptAnswer        = examdao.ExamAttemptAnswer
	ExamResult               = examdao.ExamResult
	ExamAttemptQuestionAudio = examdao.ExamAttemptQuestionAudio
	ExamBatch                = examdao.ExamBatch
	ExamBatchMockPaper       = examdao.ExamBatchMockPaper
	ExamBatchMember          = examdao.ExamBatchMember
	MockLevels               = mockdao.MockLevels
	MockExaminationPaper     = mockdao.MockExaminationPaper
	MockExaminationPart      = mockdao.MockExaminationPart
	MockExaminationSegment   = mockdao.MockExaminationSegment

	// SysMember 客户端会员（表 sys_member，原 client_user 已并入此表）
	SysMember = sysdao.SysMember
)
