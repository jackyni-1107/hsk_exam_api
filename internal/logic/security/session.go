package security

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/frame/g"

	"exam/internal/consts"
)

func sessionListKey(userType int, userId int64) string {
	return consts.SessionListKeyPrefix + userTypeTag(userType) + ":" + fmt.Sprintf("%d", userId)
}

// RegisterSession 登记新 Token，超出并发数时剔除最旧会话
func (s *sSecurity) RegisterSession(ctx context.Context, userType int, userId int64, token string, ttlSeconds int64) error {
	cfg := s.LoadSessionCfg(ctx)
	if cfg.MaxConcurrentSessions <= 0 {
		return nil
	}
	listKey := sessionListKey(userType, userId)
	_, err := g.Redis().LPush(ctx, listKey, token)
	if err != nil {
		return err
	}
	_, _ = g.Redis().Expire(ctx, listKey, ttlSeconds)

	for {
		n, err := g.Redis().LLen(ctx, listKey)
		if err != nil || n <= int64(cfg.MaxConcurrentSessions) {
			break
		}
		oldTok, err := g.Redis().RPop(ctx, listKey)
		if err != nil || oldTok.IsEmpty() {
			break
		}
		old := oldTok.String()
		tk := consts.TokenRedisKeyPrefix + userTypeTag(userType) + ":" + old
		_, _ = g.Redis().Del(ctx, tk)
	}
	return nil
}

// RemoveSessionToken 登出时从会话列表移除
func (s *sSecurity) RemoveSessionToken(ctx context.Context, userType int, userId int64, token string) {
	if token == "" || userId == 0 {
		return
	}
	listKey := sessionListKey(userType, userId)
	vals, err := g.Redis().LRange(ctx, listKey, 0, -1)
	if err != nil {
		return
	}
	for _, v := range vals {
		if v.String() == token {
			_, _ = g.Redis().LRem(ctx, listKey, 1, token)
			break
		}
	}
}

// RevokeAllUserSessions 撤销某用户全部 Token（强制下线）
func (s *sSecurity) RevokeAllUserSessions(ctx context.Context, userType int, userId int64) error {
	listKey := sessionListKey(userType, userId)
	vals, err := g.Redis().LRange(ctx, listKey, 0, -1)
	if err != nil {
		return err
	}
	tag := userTypeTag(userType)
	for _, v := range vals {
		tok := v.String()
		if tok == "" {
			continue
		}
		key := consts.TokenRedisKeyPrefix + tag + ":" + tok
		_, _ = g.Redis().Del(ctx, key)
	}
	_, _ = g.Redis().Del(ctx, listKey)
	return nil
}
