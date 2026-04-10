package security

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"exam/internal/consts"
	"exam/internal/utility"
)

// ValidatePasswordNotInHistory 新口令不能与近期历史相同
func (s *sSecurity) ValidatePasswordNotInHistory(ctx context.Context, userType int, userId int64, plainPassword string) error {
	cfg := s.LoadPasswordCfg(ctx)
	if cfg.HistoryCount <= 0 || userId <= 0 {
		return nil
	}
	rows, err := g.DB().Ctx(ctx).Model(consts.TableSysPasswordHistory).
		Where("user_type", userType).
		Where("user_id", userId).
		OrderDesc("id").
		Limit(cfg.HistoryCount).
		All()
	if err != nil {
		return err
	}
	for _, row := range rows {
		h := row["password_hash"].String()
		if utility.CheckPassword(h, plainPassword) {
			return gerror.NewCode(consts.CodePasswordReuse)
		}
	}
	return nil
}

// SavePasswordHistory 在更新密码前将旧哈希写入历史表，并裁剪超出条数
func (s *sSecurity) SavePasswordHistory(ctx context.Context, userType int, userId int64, oldHash string) error {
	cfg := s.LoadPasswordCfg(ctx)
	if cfg.HistoryCount <= 0 || oldHash == "" {
		return nil
	}
	_, err := g.DB().Ctx(ctx).Model(consts.TableSysPasswordHistory).Insert(g.Map{
		"user_type":     userType,
		"user_id":       userId,
		"password_hash": oldHash,
		"created_at":    gtime.Now(),
	})
	if err != nil {
		return err
	}
	for i := 0; i < 100; i++ {
		n, err := g.DB().Ctx(ctx).Model(consts.TableSysPasswordHistory).
			Where("user_type", userType).
			Where("user_id", userId).
			Count()
		if err != nil || n <= cfg.HistoryCount {
			break
		}
		r, err := g.DB().Exec(ctx,
			`DELETE FROM `+consts.TableSysPasswordHistory+` WHERE user_type=? AND user_id=? ORDER BY id ASC LIMIT 1`,
			userType, userId,
		)
		if err != nil {
			break
		}
		if r != nil {
			if n64, _ := r.RowsAffected(); n64 == 0 {
				break
			}
		}
	}
	return nil
}
