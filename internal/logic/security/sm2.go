package security

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"sync"

	appcfg "exam/internal/config"

	"github.com/tjfoc/gmsm/sm2"
	gmx509 "github.com/tjfoc/gmsm/x509"
)

var (
	sm2KeyOnce    sync.Once
	sm2PrivKey    *sm2.PrivateKey
	sm2PubKeyHex  string
	sm2KeyLoadErr error
)

func (s *sSecurity) DecryptLoginPassword(ctx context.Context, encrypted string) (string, error) {
	_ = ctx
	key, err := s.mustLoadSM2PrivateKey()
	if err != nil {
		return "", err
	}
	raw, err := decodeSM2Cipher(encrypted)
	if err != nil {
		return "", err
	}
	plain, err := sm2.Decrypt(key, raw, sm2.C1C3C2)
	if err == nil {
		return string(plain), nil
	}
	// 兼容部分前端库返回不带 0x04 前缀的 C1（x||y），这里补前缀重试。
	plainWithPrefix, errWithPrefix := sm2.Decrypt(key, append([]byte{0x04}, raw...), sm2.C1C3C2)
	if errWithPrefix == nil {
		return string(plainWithPrefix), nil
	}
	// 兼容部分前端库附带 0x04 前缀但服务端按无前缀解析失败的情况。
	if len(raw) > 1 && raw[0] == 0x04 {
		plain2, err2 := sm2.Decrypt(key, raw[1:], sm2.C1C3C2)
		if err2 == nil {
			return string(plain2), nil
		}
	}
	return "", err
}

func (s *sSecurity) LoginEncryptPublicKeyHex(ctx context.Context) (string, error) {
	_ = ctx
	_, err := s.mustLoadSM2PrivateKey()
	if err != nil {
		return "", err
	}
	return sm2PubKeyHex, nil
}

func (s *sSecurity) mustLoadSM2PrivateKey() (*sm2.PrivateKey, error) {
	sm2KeyOnce.Do(func() {
		cfg := appcfg.Config.SM2
		privatePem := strings.TrimSpace(cfg.PrivateKeyPem)
		if privatePem == "" {
			sm2KeyLoadErr = errors.New("sm2 private key is empty")
			return
		}
		pemBytes := []byte(privatePem)
		if !strings.Contains(privatePem, "BEGIN") {
			decoded, err := base64.StdEncoding.DecodeString(privatePem)
			if err != nil {
				sm2KeyLoadErr = err
				return
			}
			pemBytes = decoded
		}
		priv, err := gmx509.ReadPrivateKeyFromPem(pemBytes, nil)
		if err != nil {
			sm2KeyLoadErr = err
			return
		}
		sm2PrivKey = priv

		pubHex := strings.TrimSpace(cfg.PublicKeyHex)
		if pubHex == "" {
			pubHex = gmx509.WritePublicKeyToHex(&priv.PublicKey)
		}
		sm2PubKeyHex = strings.TrimSpace(pubHex)
	})
	return sm2PrivKey, sm2KeyLoadErr
}

func decodeSM2Cipher(encrypted string) ([]byte, error) {
	s := strings.TrimSpace(encrypted)
	if s == "" {
		return nil, errors.New("empty encrypted password")
	}

	if raw, err := hex.DecodeString(s); err == nil {
		return raw, nil
	}
	// 兼容 0x 前缀十六进制
	if strings.HasPrefix(strings.ToLower(s), "0x") {
		raw, err := hex.DecodeString(s[2:])
		if err == nil {
			return raw, nil
		}
	}
	raw, err := base64.StdEncoding.DecodeString(s)
	if err == nil {
		return raw, nil
	}
	return nil, errors.New("invalid encrypted password format")
}
