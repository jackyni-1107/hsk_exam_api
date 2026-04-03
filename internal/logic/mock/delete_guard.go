package mock

import (
	"github.com/gogf/gf/v2/errors/gerror"

	"exam/internal/consts"
)

// ErrMockDataMutationForbidden 禁止删除或逻辑删除 mock 真源数据（含有关联 exam 物化卷时由库 FK 兜底）。
func ErrMockDataMutationForbidden() error {
	return gerror.NewCode(consts.CodeMockDataDeleteForbidden)
}
