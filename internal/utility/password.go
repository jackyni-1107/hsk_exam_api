package utility

import (
	"crypto/rand"
	"math/big"

	"exam/internal/model/bo"

	"golang.org/x/crypto/bcrypt"
)

// CheckPassword 校验明文是否与 bcrypt 哈希匹配。
func CheckPassword(passwordHash, plainPassword string) bool {
	if passwordHash == "" || plainPassword == "" {
		return false
	}
	return bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(plainPassword)) == nil
}

const (
	passwordLowerChars   = "abcdefghjkmnpqrstuvwxyz"
	passwordUpperChars   = "ABCDEFGHJKMNPQRSTUVWXYZ"
	passwordDigitChars   = "23456789"
	passwordSpecialChars = "!@#$%^&*_-+="
	passwordDefaultMin   = 10
)

// GeneratePasswordByPolicy 按密码策略生成随机密码。
func GeneratePasswordByPolicy(cfg bo.PasswordCfg) (string, error) {
	minLen := cfg.MinLength
	if minLen < passwordDefaultMin {
		minLen = passwordDefaultMin
	}
	all := passwordLowerChars + passwordUpperChars + passwordDigitChars + passwordSpecialChars
	out := make([]byte, 0, minLen)
	if cfg.RequireLower {
		ch, err := pickRandomChar(passwordLowerChars)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}
	if cfg.RequireUpper {
		ch, err := pickRandomChar(passwordUpperChars)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}
	if cfg.RequireDigit {
		ch, err := pickRandomChar(passwordDigitChars)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}
	if cfg.RequireSpecial {
		ch, err := pickRandomChar(passwordSpecialChars)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}
	for len(out) < minLen {
		ch, err := pickRandomChar(all)
		if err != nil {
			return "", err
		}
		out = append(out, ch)
	}
	if err := shuffleBytes(out); err != nil {
		return "", err
	}
	return string(out), nil
}

func pickRandomChar(set string) (byte, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(len(set))))
	if err != nil {
		return 0, err
	}
	return set[nBig.Int64()], nil
}

func shuffleBytes(items []byte) error {
	for i := len(items) - 1; i > 0; i-- {
		jBig, err := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		if err != nil {
			return err
		}
		j := int(jBig.Int64())
		items[i], items[j] = items[j], items[i]
	}
	return nil
}
