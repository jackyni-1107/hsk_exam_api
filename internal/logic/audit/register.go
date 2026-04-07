// 向 service 注册本包实现；与 logic/exam/register.go 同规则，init 仅出现在 register.go。
package audit

import auditsvc "exam/internal/service/audit"

func init() {
	auditsvc.RegisterAudit(New())
}
