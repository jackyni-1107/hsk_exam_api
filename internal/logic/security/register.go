// 向 service 注册本包实现；与 logic/exam/register.go 同规则，init 仅出现在 register.go。
package security

import secsvc "exam/internal/service/security"

func init() {
	secsvc.RegisterSecurity(New())
}
