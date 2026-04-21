package sysuser

import (
	"context"
	"strings"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	secsvc "exam/internal/service/security"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/crypto/bcrypt"
)

func (s *sSysUser) UserList(ctx context.Context, page, size int, username string, status int) ([]sysentity.SysUser, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	model := dao.SystemUser.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
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
	var list []sysentity.SysUser
	err = model.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *sSysUser) UserRoleIds(ctx context.Context, userId int64) ([]int64, error) {
	var rows []sysentity.SysUserRole
	err := dao.SystemUserRole.Ctx(ctx).
		Where("user_id", userId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.RoleId)
	}
	return ids, nil
}

func (s *sSysUser) UserRoleIDsByUserIDs(ctx context.Context, userIDs []int64) (map[int64][]int64, error) {
	userIDs = normalizePositiveIDs(userIDs)
	if len(userIDs) == 0 {
		return map[int64][]int64{}, nil
	}
	var rows []sysentity.SysUserRole
	err := dao.SystemUserRole.Ctx(ctx).
		WhereIn("user_id", userIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		OrderAsc("id").
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	roleIDsByUser := make(map[int64][]int64, len(userIDs))
	for _, userID := range userIDs {
		roleIDsByUser[userID] = []int64{}
	}
	for _, row := range rows {
		roleIDsByUser[row.UserId] = append(roleIDsByUser[row.UserId], row.RoleId)
	}
	return roleIDsByUser, nil
}

func (s *sSysUser) UserCreate(ctx context.Context, username, password, nickname, email, mobile, creator string, status int, roleIds []int64) (int64, error) {
	username = normalizeUsername(username)
	if username == "" {
		return 0, gerror.NewCode(consts.CodeInvalidParams)
	}
	roleIds = normalizeRoleIDs(roleIds)
	if err := validateRoleIDs(ctx, roleIds); err != nil {
		return 0, err
	}
	if err := secsvc.Security().ValidatePasswordPolicy(ctx, password); err != nil {
		return 0, err
	}
	passwordHash, err := hashPassword(password)
	if err != nil {
		return 0, err
	}
	cnt, err := dao.SystemUser.Ctx(ctx).
		Wheref("username = ?", username).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeUserExists)
	}
	status = normalizeUserStatus(status)
	nickname, email, mobile = normalizeUserProfile(nickname, email, mobile)

	var id int64
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		newID, err := tx.Model(dao.SystemUser.Table()).Ctx(ctx).InsertAndGetId(sysdo.SysUser{
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
			return err
		}
		id = newID

		if len(roleIds) == 0 {
			return nil
		}
		batch := buildUserRoleBatch(id, roleIds, creator)
		_, err = tx.Model(dao.SystemUserRole.Table()).Ctx(ctx).Data(batch).Insert()
		return err
	})
	if err != nil {
		return 0, err
	}

	var after sysentity.SysUser
	if err := dao.SystemUser.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SystemUser.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sSysUser) UserUpdate(ctx context.Context, id int64, password, nickname, email, mobile, updater string, status int) error {
	before, err := loadUserByID(ctx, id)
	if err != nil {
		return err
	}
	nickname, email, mobile = normalizeUserProfile(nickname, email, mobile)
	data := sysdo.SysUser{
		Nickname: nickname,
		Email:    email,
		Mobile:   mobile,
		Updater:  updater,
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
		passwordHash, err := preparePasswordChange(ctx, before.Id, before.Password, password)
		if err != nil {
			return err
		}
		data.Password = passwordHash
		data.PasswordChangedAt = gtime.Now()
		shouldRevokeSessions = true
		passwordChanged = true
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemUser.Table()).Ctx(ctx).Where("id", id).Data(data).Update(); err != nil {
			return err
		}
		if passwordChanged {
			return secsvc.Security().SavePasswordHistory(ctx, consts.UserTypeAdmin, id, before.Password)
		}
		return nil
	})
	if err != nil {
		return err
	}
	if shouldRevokeSessions {
		bestEffortRevokeAdminSessions(ctx, id)
	}
	var after sysentity.SysUser
	if err := dao.SystemUser.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemUser.Table(), id, before, &after)
	}
	return nil
}

func (s *sSysUser) UserDelete(ctx context.Context, id int64, updater string) error {
	if id == consts.SuperAdminUserId {
		return gerror.NewCode(consts.CodeCannotDeleteSuperAdmin)
	}
	before, err := loadUserByID(ctx, id)
	if err != nil {
		return err
	}
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemUser.Table()).Ctx(ctx).Where("id", id).Data(sysdo.SysUser{
			DeleteFlag: consts.DeleteFlagDeleted,
			Updater:    updater,
		}).Update(); err != nil {
			return err
		}
		_, err := tx.Model(dao.SystemUserRole.Table()).Ctx(ctx).
			Where("user_id", id).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysUserRole{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    updater,
			}).Update()
		return err
	})
	if err != nil {
		return err
	}
	bestEffortClearUserPermissionCache(ctx, id)
	bestEffortRevokeAdminSessions(ctx, id)
	var after sysentity.SysUser
	if err := dao.SystemUser.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemUser.Table(), id, before, &after)
	}
	return nil
}

func (s *sSysUser) UserRoleAssign(ctx context.Context, userId int64, roleIds []int64, creator string) error {
	if _, err := loadUserByID(ctx, userId); err != nil {
		return err
	}
	roleIds = normalizeRoleIDs(roleIds)
	if err := validateRoleIDs(ctx, roleIds); err != nil {
		return err
	}
	beforeIDs, err := s.UserRoleIds(ctx, userId)
	if err != nil {
		return err
	}
	beforeStr := utility.JoinSortedInt64IDs(beforeIDs)
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := tx.Model(dao.SystemUserRole.Table()).Ctx(ctx).
			Where("user_id", userId).
			Where("delete_flag", consts.DeleteFlagNotDeleted).
			Data(sysdo.SysUserRole{
				DeleteFlag: consts.DeleteFlagDeleted,
				Updater:    creator,
			}).Update(); err != nil {
			return err
		}
		if len(roleIds) == 0 {
			return nil
		}
		batch := buildUserRoleBatch(userId, roleIds, creator)
		_, err := tx.Model(dao.SystemUserRole.Table()).Ctx(ctx).Data(batch).Insert()
		return err
	})
	if err != nil {
		return err
	}
	bestEffortClearUserPermissionCache(ctx, userId)
	afterStr := utility.JoinSortedInt64IDs(roleIds)
	auditutil.RecordMapDiff(ctx, dao.SystemUserRole.Table(), userId,
		map[string]interface{}{"role_ids": beforeStr},
		map[string]interface{}{"role_ids": afterStr})
	return nil
}

func (s *sSysUser) FindByUsername(ctx context.Context, username string) (*sysentity.SysUser, error) {
	username = normalizeUsername(username)
	if username == "" {
		return nil, gerror.NewCode(consts.CodeInvalidParams)
	}
	var u sysentity.SysUser
	err := dao.SystemUser.Ctx(ctx).
		Where("username", username).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&u)
	if err != nil {
		return nil, err
	}
	if u.Id == 0 {
		return nil, nil
	}
	return &u, nil
}

func normalizeUsername(username string) string {
	return strings.TrimSpace(username)
}

func normalizeUserProfile(nickname, email, mobile string) (string, string, string) {
	return strings.TrimSpace(nickname), strings.TrimSpace(email), strings.TrimSpace(mobile)
}

func normalizeUserStatus(status int) int {
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		return status
	}
	return consts.StatusNormal
}

func normalizeRoleIDs(roleIDs []int64) []int64 {
	return normalizePositiveIDs(roleIDs)
}

func normalizePositiveIDs(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

func validateRoleIDs(ctx context.Context, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return nil
	}
	cnt, err := dao.SystemRole.Ctx(ctx).
		WhereIn("id", roleIDs).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Where("status", consts.StatusNormal).
		Count()
	if err != nil {
		return err
	}
	if cnt != len(roleIDs) {
		return gerror.NewCode(consts.CodeRoleNotFound)
	}
	return nil
}

func loadUserByID(ctx context.Context, userID int64) (*sysentity.SysUser, error) {
	var user sysentity.SysUser
	if err := dao.SystemUser.Ctx(ctx).
		Where("id", userID).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&user); err != nil {
		return nil, err
	}
	if user.Id == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	return &user, nil
}

func preparePasswordChange(ctx context.Context, userID int64, currentHash, newPassword string) (string, error) {
	if err := secsvc.Security().ValidatePasswordPolicy(ctx, newPassword); err != nil {
		return "", err
	}
	if utility.CheckPassword(currentHash, newPassword) {
		return "", gerror.NewCode(consts.CodePasswordReuse)
	}
	if err := secsvc.Security().ValidatePasswordNotInHistory(ctx, consts.UserTypeAdmin, userID, newPassword); err != nil {
		return "", err
	}
	return hashPassword(newPassword)
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func buildUserRoleBatch(userID int64, roleIDs []int64, actor string) []sysdo.SysUserRole {
	batch := make([]sysdo.SysUserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		batch = append(batch, sysdo.SysUserRole{
			UserId:  userID,
			RoleId:  roleID,
			Creator: actor,
			Updater: actor,
		})
	}
	return batch
}

func bestEffortClearUserPermissionCache(ctx context.Context, userID int64) {
	if userID <= 0 {
		return
	}
	if _, err := g.Redis().Del(ctx, consts.PermCacheKeyPrefix+gconv.String(userID)); err != nil {
		g.Log().Warningf(ctx, "clear user permission cache failed user_id=%d: %v", userID, err)
	}
}

func bestEffortRevokeAdminSessions(ctx context.Context, userID int64) {
	if userID <= 0 {
		return
	}
	if err := secsvc.Security().RevokeAllUserSessions(ctx, consts.UserTypeAdmin, userID); err != nil {
		g.Log().Warningf(ctx, "revoke admin sessions failed user_id=%d: %v", userID, err)
	}
}
