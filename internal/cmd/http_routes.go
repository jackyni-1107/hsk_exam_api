package cmd

import (
	"github.com/gogf/gf/v2/net/ghttp"

	appcfg "exam/internal/config"
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
	"exam/internal/utility/openapi"
)

func registerHTTPRoutes(s *ghttp.Server) {
	apiPrefix := appcfg.Config.ApiPrefix
	registerHealthRoutes(s)
	registerClientRoutes(s, apiPrefix)
	registerAdminRoutes(s, apiPrefix)
	openapi.RegisterSplitEndpoints(s, apiPrefix)
}

func registerHealthRoutes(s *ghttp.Server) {
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.GET("/health", health.Liveness)
		group.GET("/ready", health.Readiness)
	})
}

func registerClientRoutes(s *ghttp.Server, apiPrefix string) {
	clientAPI := appcfg.JoinHTTPPath(apiPrefix, "client")
	clientExamMedia := appcfg.JoinHTTPPath(apiPrefix, "client/exam/media")

	s.Group(clientAPI, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.ClientPublicChain()...)
		group.Bind(clientAuth.NewV1())
	})
	s.Group(clientAPI, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.ClientProtectedChain()...)
		group.Bind(clientMe.NewV1(), clientExam.NewV1())
	})
	s.Group(clientExamMedia, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.ClientMediaChain()...)
		group.GET("/hls/{ticket}.m3u8", clientExam.ServeHlsM3U8)
	})
}

func registerAdminRoutes(s *ghttp.Server, apiPrefix string) {
	adminAPI := appcfg.JoinHTTPPath(apiPrefix, "admin")

	s.Group(adminAPI, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.AdminPublicChain()...)
		group.Bind(adminAuth.NewV1())
	})
	s.Group(adminAPI, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.AdminProtectedChain()...)
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
	s.Group(adminAPI, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.AdminDownloadChain()...)
		group.GET("/file/{id}/download", adminFile.ServeDownload)
		group.GET("/member/import-template", adminMember.ServeMemberImportTemplate)
		group.GET("/exam/batch/member-import-template", adminExam.ServeBatchMemberImportTemplate)
	})
}
