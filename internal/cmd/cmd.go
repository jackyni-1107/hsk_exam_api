package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"

	"exam/internal/consts"
	adminAuditLog "exam/internal/controller/admin/audit_log"
	adminAuth "exam/internal/controller/admin/auth"
	adminConfig "exam/internal/controller/admin/config"
	adminExam "exam/internal/controller/admin/exam"
	adminExceptionLog "exam/internal/controller/admin/exception_log"
	adminFile "exam/internal/controller/admin/file"
	adminLoginLog "exam/internal/controller/admin/login_log"
	adminMe "exam/internal/controller/admin/me"
	adminMember "exam/internal/controller/admin/member"
	adminMenu "exam/internal/controller/admin/menu"
	adminMock "exam/internal/controller/admin/mock"
	adminNotification "exam/internal/controller/admin/notification"
	adminRole "exam/internal/controller/admin/role"
	adminSecurityEventLog "exam/internal/controller/admin/security_event_log"
	adminTask "exam/internal/controller/admin/task"
	adminUser "exam/internal/controller/admin/user"
	clientAuth "exam/internal/controller/client/auth"
	clientExam "exam/internal/controller/client/exam"
	clientMe "exam/internal/controller/client/me"
	"exam/internal/controller/health"
	"exam/internal/middleware"
	"exam/internal/openapi"
	"exam/internal/task"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			InitAll(ctx)

			s := g.Server()
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.GET("/health", health.Liveness)
				group.GET("/ready", health.Readiness)
			})
			// 客户端认证：login / captcha / logout 不经过 Auth，无需 Token。
			s.Group("/api/client", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Trace, middleware.Response, middleware.HandlerResponseI18n)
				group.Bind(clientAuth.NewV1())
			})
			s.Group("/api/client", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Trace, middleware.Response, middleware.HandlerResponseI18n, middleware.Auth(consts.UserTypeClient))
				group.Bind(clientMe.NewV1(), clientExam.NewV1())
			})
			s.Group("/api/client/exam/media", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Trace, middleware.Response)
				group.GET("/hls/{ticket}.m3u8", clientExam.ServeHlsM3U8)
			})
			s.Group("/api/admin", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.Trace, middleware.Response, middleware.HandlerResponseI18n)
				group.Bind(adminAuth.NewV1())
			})
			s.Group("/api/admin", func(group *ghttp.RouterGroup) {
				group.Middleware(
					middleware.Trace,
					middleware.Response,
					middleware.HandlerResponseI18n,
					middleware.Auth(consts.UserTypeAdmin),
					middleware.RBACFromPath,
					middleware.Audit,
				)
				group.Bind(
					adminMe.NewV1(),
					adminUser.NewV1(),
					adminRole.NewV1(),
					adminMenu.NewV1(),
					adminMember.NewV1(),
					adminNotification.NewV1(),
					adminFile.NewV1(),
					adminConfig.NewV1(),
					adminTask.NewV1(),
					adminAuditLog.NewV1(),
					adminLoginLog.NewV1(),
					adminExceptionLog.NewV1(),
					adminSecurityEventLog.NewV1(),
					adminExam.NewV1(),
					adminMock.NewV1(),
				)
			})
			openapi.RegisterSplitEndpoints(s)
			task.StartScheduler(ctx)
			s.Run()
			return nil
		},
	}
)
