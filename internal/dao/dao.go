package dao

import (
	daoexam "exam/internal/dao/exam"
	daomock "exam/internal/dao/mock"
	daosys "exam/internal/dao/sys"
)

// 管理端 RBAC 等使用的 System* 命名（底层为 sys_* DAO）。
var (
	SystemUser     = daosys.SysUser
	SystemUserRole = daosys.SysUserRole
	SystemRoleMenu = daosys.SysRoleMenu
	SystemMenu     = daosys.SysMenu
	SystemRole     = daosys.SysRole

	SystemDictType = daosys.SysDictType
	SystemDictData = daosys.SysDictData
	SystemConfig   = daosys.SysConfig

	SystemOperationAuditLog = daosys.SysOperationAuditLog
)

// Sys* 表（任务、文件、通知等）
var (
	SysTask                      = daosys.SysTask
	SysTaskLog                   = daosys.SysTaskLog
	SysFileStorage               = daosys.SysFileStorage
	SysFileStorageConfig         = daosys.SysFileStorageConfig
	SysNotificationChannelConfig = daosys.SysNotificationChannelConfig
	SysNotificationTemplate      = daosys.SysNotificationTemplate
	SysNotificationLog           = daosys.SysNotificationLog
)

// 考试与模拟卷
var (
	ExamPaper                = daoexam.ExamPaper
	ExamSection              = daoexam.ExamSection
	ExamQuestionBlock        = daoexam.ExamQuestionBlock
	ExamQuestion             = daoexam.ExamQuestion
	ExamOption               = daoexam.ExamOption
	ExamAttempt              = daoexam.ExamAttempt
	ExamAttemptAnswer        = daoexam.ExamAttemptAnswer
	ExamResult               = daoexam.ExamResult
	ExamAttemptQuestionAudio = daoexam.ExamAttemptQuestionAudio
	ExamBatch                = daoexam.ExamBatch
	ExamBatchMockLevel       = daoexam.ExamBatchMockLevel
	ExamBatchMember          = daoexam.ExamBatchMember
	MockLevels               = daomock.MockLevels
	MockExaminationPaper     = daomock.MockExaminationPaper
	MockExaminationPart      = daomock.MockExaminationPart
	MockExaminationSegment   = daomock.MockExaminationSegment

	// SysMember 客户端会员（表 sys_member，原 client_user 已并入此表）
	SysMember = daosys.SysMember
)
