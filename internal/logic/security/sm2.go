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
	if plain, ok := tryDecryptSM2Compat(key, raw); ok {
		return string(plain), nil
	}
	return "", errors.New("sm2 decrypt failed for all supported modes")
}

func tryDecryptSM2Compat(key *sm2.PrivateKey, raw []byte) ([]byte, bool) {
	candidates := make([][]byte, 0, 3)
	candidates = append(candidates, raw)
	// 兼容部分前端库返回不带 0x04 前缀的 C1（x||y），这里补前缀重试。
	candidates = append(candidates, append([]byte{0x04}, raw...))
	// 兼容部分前端库附带 0x04 前缀但服务端按无前缀解析失败的情况。
	if len(raw) > 1 && raw[0] == 0x04 {
		candidates = append(candidates, raw[1:])
	}

	modeList := []int{sm2.C1C3C2, sm2.C1C2C3}
	for _, m := range modeList {
		for _, c := range candidates {
			if p, ok := safeDecryptByMode(key, c, m); ok {
				return p, true
			}
		}
	}
	return nil, false
}

func safeDecryptByMode(key *sm2.PrivateKey, raw []byte, mode int) (plain []byte, ok bool) {
	defer func() {
		if recover() != nil {
			plain = nil
			ok = false
		}
	}()
	p, err := sm2.Decrypt(key, raw, mode)
	if err != nil {
		return nil, false
	}
	return p, true
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
