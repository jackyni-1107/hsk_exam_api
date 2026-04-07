package config

import "exam/internal/model/bo"

// Config 进程内缓存的配置快照（由 cmd.InitAll 从 g.Cfg 填充）。
var Config struct {
	Exam     bo.ExamCfg
	Login    bo.LoginCfg
	Password bo.PasswordCfg
	Session  bo.SessionCfg
	MFA      bo.MFACfg
}
