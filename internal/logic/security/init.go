// 向 service 注册本包实现；与 logic/exam/init.go 同规则，init 仅出现在 init.go。
package security

import "exam/internal/service/security"

type sSecurity struct{}

func init() {
	security.RegisterSecurity(New())
}

func New() *sSecurity {
	return &sSecurity{}
}
