package security

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"

	"exam/internal/consts"
	"exam/internal/model/bo"
)

// CreateCaptcha 生成验证码并存 Redis
func (s *sSecurity) CreateCaptcha(ctx context.Context) (*bo.CaptchaChallenge, error) {
	a := rand.Intn(20) + 1
	b := rand.Intn(20) + 1
	answer := a + b
	id := guid.S()
	key := consts.CaptchaKeyPrefix + id
	if err := g.Redis().SetEX(ctx, key, fmt.Sprintf("%d", answer), consts.CaptchaTTLSeconds); err != nil {
		return nil, err
	}
	return &bo.CaptchaChallenge{
		CaptchaId: id,
		Question:  fmt.Sprintf("%d + %d = ?", a, b),
	}, nil
}

// VerifyCaptcha 校验答案（一次性）
func (s *sSecurity) VerifyCaptcha(ctx context.Context, captchaId, answer string) bool {
	if captchaId == "" || answer == "" {
		return false
	}
	key := consts.CaptchaKeyPrefix + captchaId
	val, err := g.Redis().Get(ctx, key)
	if err != nil || val.IsEmpty() {
		return false
	}
	ok := val.String() == answer
	_, _ = g.Redis().Del(ctx, key)
	return ok
}
