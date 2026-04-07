package security

// sSecurity 实现 service/security.ISecurity，方法分布在各业务文件中。
type sSecurity struct{}

func New() *sSecurity {
	return &sSecurity{}
}
