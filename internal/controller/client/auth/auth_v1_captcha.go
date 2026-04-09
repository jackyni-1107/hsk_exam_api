package auth

import (
	"context"

	v1 "exam/api/client/auth/v1"
	"exam/internal/pkg/captchaimg"
	secsvc "exam/internal/service/security"
)

func (c *ControllerV1) Captcha(ctx context.Context, req *v1.CaptchaReq) (res *v1.CaptchaRes, err error) {
	ch, err := secsvc.Security().CreateCaptcha(ctx)
	if err != nil {
		return nil, err
	}
	imgB64, err := captchaimg.QuestionToPNGBase64(ch.Question)
	if err != nil {
		return nil, err
	}
	return &v1.CaptchaRes{CaptchaId: ch.CaptchaId, Question: ch.Question, QuestionImage: imgB64}, nil
}
