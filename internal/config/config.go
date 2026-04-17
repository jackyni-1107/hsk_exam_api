package config

import "exam/internal/model/bo"

// Config 进程内缓存的配置快照（由 cmd.InitAll 从 g.Cfg 填充）。
var Config struct {
	// ApiPrefix HTTP 路由前缀，如 "/api"。开发环境直连时常用；测试/生产若由 nginx 剥前缀后转发，可设为 ""。
	ApiPrefix string
	Exam      bo.ExamCfg
	Login     bo.LoginCfg
	Password  bo.PasswordCfg
	Session   bo.SessionCfg
	SM2       bo.SM2Cfg
}
