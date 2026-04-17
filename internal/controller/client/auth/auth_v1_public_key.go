package auth

import (
	"context"

	v1 "exam/api/client/auth/v1"
	"exam/internal/consts"
	secsvc "exam/internal/service/security"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (c *ControllerV1) PublicKey(ctx context.Context, req *v1.PublicKeyReq) (res *v1.PublicKeyRes, err error) {
	_ = req
	publicKeyHex, err := secsvc.Security().LoginEncryptPublicKeyHex(ctx)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeLoginFailed)
	}
	return &v1.PublicKeyRes{
		PublicKeyHex: publicKeyHex,
		Algorithm:    "sm2",
		CipherMode:   "c1c3c2",
	}, nil
}
