package storage

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/frame/g"
)

// normalizePresignSigVersion 将配置中的签名版本归一为 v2 / v3 / v4。
// v3 与 v4 使用相同实现（AWS S3 REST 预签名无独立 SigV3）。
// 无法识别的值记日志并回退为 v4。
func normalizePresignSigVersion(ctx context.Context, raw string) string {
	s := strings.TrimSpace(strings.ToLower(raw))
	switch s {
	case "", "v4", "4", "sigv4", "s3v4":
		return "v4"
	case "v3", "3", "sigv3", "s3v3":
		return "v3"
	case "v2", "2", "sigv2", "s3", "s3v2":
		return "v2"
	default:
		if s != "" {
			g.Log().Warningf(ctx, "storage: unknown presign_signature_version %q, using v4", raw)
		}
		return "v4"
	}
}

func applyPublicBaseURL(raw string, publicBase string) (string, error) {
	if strings.TrimSpace(publicBase) == "" {
		return raw, nil
	}
	return rewriteURLOrigin(raw, strings.TrimRight(strings.TrimSpace(publicBase), "/"))
}
