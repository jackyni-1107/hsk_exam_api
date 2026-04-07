// 向 service 注册本包实现；与 logic/audit/register.go 同规则，init 仅出现在 register.go。
package exam

import examsvc "exam/internal/service/exam"

func init() {
	examsvc.RegisterExam(new(sExam))
}
