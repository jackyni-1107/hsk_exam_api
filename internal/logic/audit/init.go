// 向 service 注册本包实现；与 logic/exam/init.go 同规则，init 仅出现在 init.go。
package audit

import "exam/internal/service/audit"

type sAudit struct{}

func New() *sAudit {
	return &sAudit{}
}

func init() {
	audit.RegisterAudit(New())
}
