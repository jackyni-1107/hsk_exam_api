package security

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"exam/internal/consts"
	"exam/internal/model/bo"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
)

func tokenKey(userType int, token string) string {
	return consts.TokenRedisKeyPrefix + userTypeTag(userType) + ":" + token
}

func decodeTokenPayload(raw []byte) (*bo.TokenPayload, error) {
	var payload bo.TokenPayload
	if err := json.Unmarshal(raw, &payload); err != nil {
		return nil, err
	}
	if payload.UserId == 0 {
		return nil, fmt.Errorf("invalid token payload")
	}
	return &payload, nil
}

func (s *sSecurity) IssueToken(ctx context.Context, userType int, userId int64, username string) (string, error) {
	token := guid.S()
	ttl := s.TokenTTLSeconds(ctx)
	if ttl <= 0 {
		ttl = consts.DefaultTokenTTLFallbackSeconds
	}
	payload, err := json.Marshal(bo.TokenPayload{
		UserId:   userId,
		Username: username,
	})
	if err != nil {
		return "", err
	}
	key := tokenKey(userType, token)
	if err := g.Redis().SetEX(ctx, key, string(payload), ttl); err != nil {
		return "", err
	}
	if err := s.RegisterSession(ctx, userType, userId, token, ttl); err != nil {
		_, _ = g.Redis().Del(ctx, key)
		return "", err
	}
	return token, nil
}

func (s *sSecurity) LoadTokenPayload(ctx context.Context, userType int, token string) (*bo.TokenPayload, error) {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil, nil
	}
	val, err := g.Redis().Get(ctx, tokenKey(userType, token))
	if err != nil {
		return nil, err
	}
	if val.IsEmpty() {
		return nil, nil
	}
	return decodeTokenPayload(val.Bytes())
}

func (s *sSecurity) RevokeToken(ctx context.Context, userType int, userId int64, token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return nil
	}
	if userId != 0 {
		s.RemoveSessionToken(ctx, userType, userId, token)
	}
	_, err := g.Redis().Del(ctx, tokenKey(userType, token))
	return err
}
