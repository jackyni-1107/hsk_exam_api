package security

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"

	"exam/internal/utility"

	"github.com/tjfoc/gmsm/sm2"
)

func (s *sSecurity) EncryptMemberPassword(ctx context.Context, plain string) (string, error) {
	_ = ctx
	_, err := s.mustLoadSM2PrivateKey()
	if err != nil {
		return "", err
	}
	cipher, err := sm2.Encrypt(&sm2PrivKey.PublicKey, []byte(plain), rand.Reader, sm2.C1C3C2)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(cipher), nil
}

func (s *sSecurity) VerifyMemberPassword(ctx context.Context, encrypted string, plain string) (bool, error) {
	_ = ctx
	encrypted = strings.TrimSpace(encrypted)
	if encrypted == "" {
		return false, nil
	}

	// 兼容历史 bcrypt 存量，避免切换期间老账号无法登录。
	if strings.HasPrefix(encrypted, "$2a$") || strings.HasPrefix(encrypted, "$2b$") || strings.HasPrefix(encrypted, "$2y$") {
		return utility.CheckPassword(encrypted, plain), nil
	}
	// SM2 密文（hex）正常长度应明显大于 bcrypt；过短通常是数据库字段截断导致。
	if len(encrypted) < 180 {
		return false, errors.New("member password cipher too short, possible truncated data")
	}

	raw, err := decodeSM2Cipher(encrypted)
	if err != nil {
		return false, err
	}
	key, err := s.mustLoadSM2PrivateKey()
	if err != nil {
		return false, err
	}
	dec, ok := tryDecryptSM2Compat(key, raw)
	if ok {
		return string(dec) == plain, nil
	}
	return false, nil
}
