package do

import (
	doexam "exam/internal/model/do/exam"
	dosys "exam/internal/model/do/sys"
)

type (
	ExamAttempt       = doexam.ExamAttempt
	ExamAttemptAnswer = doexam.ExamAttemptAnswer
	ExamPaper         = doexam.ExamPaper
)

type (
	SysTask                      = dosys.SysTask
	SysTaskLog                   = dosys.SysTaskLog
	SysFileStorage               = dosys.SysFileStorage
	SysFileStorageConfig         = dosys.SysFileStorageConfig
	SysNotificationChannelConfig = dosys.SysNotificationChannelConfig
	SysNotificationTemplate      = dosys.SysNotificationTemplate
	SysConfig                    = dosys.SysConfig
	SysDictType                  = dosys.SysDictType
	SysDictData                  = dosys.SysDictData
)

// 与历史命名 System* 对齐
type (
	SystemConfig   = dosys.SysConfig
	SystemDictType = dosys.SysDictType
	SystemDictData = dosys.SysDictData
)
