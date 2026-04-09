package cmd

import (
	"context"
	"strings"

	appcfg "exam/internal/config"
	"exam/internal/model/bo"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

// InitAll 初始化所有基础设施（启动时调用）
func InitAll(ctx context.Context) {
	initConfig(ctx)
}

// 初始化配置
func initConfig(ctx context.Context) {
	configAdapter, err := gcfg.NewAdapterFile("config.yaml")
	if err != nil {
		panic(err)
	}
	g.Cfg().SetAdapter(configAdapter)

	if err = g.Cfg().MustGet(ctx, "security.login").Scan(&appcfg.Config.Login); err != nil {
		panic(err)
	}
	if err = g.Cfg().MustGet(ctx, "security.password").Scan(&appcfg.Config.Password); err != nil {
		panic(err)
	}
	if err = g.Cfg().MustGet(ctx, "security.session").Scan(&appcfg.Config.Session); err != nil {
		panic(err)
	}
	if err = g.Cfg().MustGet(ctx, "security.mfa").Scan(&appcfg.Config.MFA); err != nil {
		panic(err)
	}
	var examCfg bo.ExamCfg
	if err = g.Cfg().MustGet(ctx, "exam").Scan(&examCfg); err != nil {
		panic(err)
	}
	appcfg.Config.Exam = examCfg

	initApiPrefix(ctx)
}

// initApiPrefix 读取 server.apiPrefix；未配置时默认 "/api" 与历史行为一致；显式配置为空字符串表示无前缀。
func initApiPrefix(ctx context.Context) {
	v, err := g.Cfg().Get(ctx, "server.apiPrefix")
	if err != nil || v.IsNil() {
		appcfg.Config.ApiPrefix = "/api"
		return
	}
	appcfg.Config.ApiPrefix = strings.TrimSuffix(strings.TrimSpace(v.String()), "/")
}
