package member

import (
	"context"
	"strings"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	"exam/internal/middleware"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	secsvc "exam/internal/service/security"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

func (s *sMember) MemberList(ctx context.Context, page, size int, username string, status int) ([]sysentity.SysMember, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	username = normalizeMemberUsername(username)
	model := dao.SysMember.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if username != "" {
		model = model.WhereLike("username", "%"+username+"%")
	}
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		model = model.Where("status", status)
	}
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []sysentity.SysMember
	if err = model.Page(page, size).OrderDesc("id").Scan(&list); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *sMember) MemberCreate(ctx context.Context, username, password, nickname, email, mobile, creator string, status int) (int64, error) {
	username = normalizeMemberUsername(username)
	if username == "" {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	//if err := secsvc.Security().ValidatePasswordPolicy(ctx, password); err != nil {
	//	return 0, err
	//}
	passwordHash, err := encryptMemberPassword(ctx, password)
	if err != nil {
		return 0, err
	}
	cnt, err := dao.SysMember.Ctx(ctx).
		Where("username", username).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeMemberExists)
	}

	status = normalizeMemberStatus(status)
	nickname, email, mobile = normalizeMemberProfile(nickname, email, mobile)
	creator = sanitizeActorByContext(ctx, creator)
	id, err := dao.SysMember.Ctx(ctx).InsertAndGetId(sysdo.SysMember{
		Username:          username,
		Password:          passwordHash,
		PasswordChangedAt: gtime.Now(),
		Nickname:          nickname,
		Email:             email,
		Mobile:            mobile,
		Status:            status,
		Creator:           creator,
		Updater:           creator,
	})
	if err != nil {
		return 0, err
	}
	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sMember) MemberUpdate(ctx context.Context, id int64, password, nickname, email, mobile, updater string, status int) error {
	before, err := loadMemberByID(ctx, id)
	if err != nil {
		return err
	}

	nickname, email, mobile = normalizeMemberProfile(nickname, email, mobile)
	data := sysdo.SysMember{
		Nickname: nickname,
		Email:    email,
		Mobile:   mobile,
	}
	if actor := sanitizeActorByContext(ctx, updater); actor != "" {
		data.Updater = actor
	}
	shouldRevokeSessions := false
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		data.Status = status
		if before.Status != status && status == consts.StatusDisabled {
			shouldRevokeSessions = true
		}
	}

	passwordChanged := false
	if password != "" {
		passwordHash, err := encryptMemberPassword(ctx, password)
		if err != nil {
			return err
		}
		data.Password = passwordHash
		data.PasswordChangedAt = gtime.Now()
		shouldRevokeSessions = true
		passwordChanged = true
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SysMember.Table()).Ctx(ctx).Where("id", id).Data(data).Update(); err != nil {
			return err
		}
		if passwordChanged {
			return secsvc.Security().SavePasswordHistory(ctx, consts.UserTypeClient, id, before.Password)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if shouldRevokeSessions {
		bestEffortRevokeClientSessions(ctx, id)
	}

	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, before, &after)
	}
	return nil
}

func (s *sMember) MemberUpdatePwd(ctx context.Context, id int64, password string) error {
	before, err := loadMemberByID(ctx, id)
	if err != nil {
		return err
	}

	data := sysdo.SysMember{
		Password: password,
	}
	shouldRevokeSessions := false
	passwordChanged := false
	if password != "" {
		passwordHash, err := encryptMemberPassword(ctx, password)
		if err != nil {
			return err
		}
		data.Password = passwordHash
		data.PasswordChangedAt = gtime.Now()
		shouldRevokeSessions = true
		passwordChanged = true
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SysMember.Table()).Ctx(ctx).Where("id", id).Data(data).Update(); err != nil {
			return err
		}
		if passwordChanged {
			return secsvc.Security().SavePasswordHistory(ctx, consts.UserTypeClient, id, before.Password)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if shouldRevokeSessions {
		bestEffortRevokeClientSessions(ctx, id)
	}

	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, before, &after)
	}
	return nil
}

func (s *sMember) MemberDelete(ctx context.Context, id int64, updater string) error {
	before, err := loadMemberByID(ctx, id)
	if err != nil {
		return err
	}
	_, err = dao.SysMember.Ctx(ctx).Where("id", id).Data(sysdo.SysMember{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return err
	}
	bestEffortRevokeClientSessions(ctx, id)

	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, before, &after)
	}
	return nil
}

func (s *sMember) MemberProfile(ctx context.Context, memberId int64) (*sysentity.SysMember, error) {
	return loadMemberByID(ctx, memberId)
}

func (s *sMember) FindByUsername(ctx context.Context, username string) (*sysentity.SysMember, error) {
	username = normalizeMemberUsername(username)
	if username == "" {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	var member sysentity.SysMember
	err := dao.SysMember.Ctx(ctx).
		Where("username", username).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&member)
	if err != nil {
		return nil, err
	}
	if member.Id == 0 {
		return nil, nil
	}
	return &member, nil
}

func (s *sMember) FindByEmail(ctx context.Context, email string) (*sysentity.SysMember, error) {
	var member sysentity.SysMember
	err := dao.SysMember.Ctx(ctx).
		Where("email", email).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&member)
	if err != nil {
		return nil, err
	}
	if member.Id == 0 {
		return nil, nil
	}
	return &member, nil
}

func normalizeMemberUsername(username string) string {
	return strings.TrimSpace(username)
}

func normalizeMemberProfile(nickname, email, mobile string) (string, string, string) {
	return strings.TrimSpace(nickname), strings.TrimSpace(email), strings.TrimSpace(mobile)
}

func normalizeMemberStatus(status int) int {
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		return status
	}
	return consts.StatusNormal
}

func loadMemberByID(ctx context.Context, memberID int64) (*sysentity.SysMember, error) {
	var member sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).
		Where("id", memberID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&member); err != nil {
		return nil, err
	}
	if member.Id == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	return &member, nil
}

func encryptMemberPassword(ctx context.Context, password string) (string, error) {
	hash, err := secsvc.Security().EncryptMemberPassword(ctx, password)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func bestEffortRevokeClientSessions(ctx context.Context, memberID int64) {
	if memberID <= 0 {
		return
	}
	if err := secsvc.Security().RevokeAllUserSessions(ctx, consts.UserTypeClient, memberID); err != nil {
		g.Log().Warningf(ctx, "revoke client sessions failed member_id=%d: %v", memberID, err)
	}
}

func sanitizeActorByContext(ctx context.Context, actor string) string {
	actor = strings.TrimSpace(actor)
	if actor == "" {
		return ""
	}
	if d := middleware.GetCtxData(ctx); d != nil && strings.TrimSpace(d.Username) == actor {
		return actor
	}
	return ""
}
