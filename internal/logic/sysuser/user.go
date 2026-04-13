package sysuser

import (
	"context"

	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
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

func (s *sSysUser) UserCreate(ctx context.Context, username, password, nickname, email, mobile, creator string, status int, roleIds []int64) (int64, error) {
	cnt, err := dao.SystemUser.Ctx(ctx).
		Where("username", username).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeUserExists)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	if status != consts.StatusNormal && status != consts.StatusDisabled {
		status = consts.StatusNormal
	}

	id, err := dao.SystemUser.Ctx(ctx).InsertAndGetId(sysdo.SysUser{
		Username: username,
		Password: string(hash),
		Nickname: nickname,
		Email:    email,
		Mobile:   mobile,
		Status:   status,
		Creator:  creator,
		Updater:  creator,
	})
	if err != nil {
		return 0, err
	}

	if len(roleIds) > 0 {
		batch := make([]sysdo.SysUserRole, 0, len(roleIds))
		for _, rid := range roleIds {
			batch = append(batch, sysdo.SysUserRole{
				UserId:  id,
				RoleId:  rid,
				Creator: creator,
				Updater: creator,
			})
		}
		_, err = dao.SystemUserRole.Ctx(ctx).Insert(batch)
		if err != nil {
			return 0, err
		}
	}
	return id, nil
}

func (s *sSysUser) UserUpdate(ctx context.Context, id int64, password, nickname, email, mobile, updater string, status int) error {
	data := sysdo.SysUser{
		Nickname: nickname,
		Email:    email,
		Mobile:   mobile,
		Updater:  updater,
	}
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		data.Status = status
	}
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		data.Password = string(hash)
	}
	_, err := dao.SystemUser.Ctx(ctx).Where("id", id).Data(data).Update()
	return err
}

func (s *sSysUser) UserDelete(ctx context.Context, id int64, updater string) error {
	if id == consts.SuperAdminUserId {
		return gerror.NewCode(consts.CodeCannotDeleteSuperAdmin)
	}
	_, err := dao.SystemUser.Ctx(ctx).Where("id", id).Data(sysdo.SysUser{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	return err
}

func (s *sSysUser) UserRoleAssign(ctx context.Context, userId int64, roleIds []int64, creator string) error {
	_, err := dao.SystemUserRole.Ctx(ctx).
		Where("user_id", userId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(sysdo.SysUserRole{
			DeleteFlag: consts.DeleteFlagDeleted,
			Updater:    creator,
		}).Update()
	if err != nil {
		return err
	}

	if len(roleIds) > 0 {
		batch := make([]sysdo.SysUserRole, 0, len(roleIds))
		for _, rid := range roleIds {
			batch = append(batch, sysdo.SysUserRole{
				UserId:  userId,
				RoleId:  rid,
				Creator: creator,
				Updater: creator,
			})
		}
		_, err = dao.SystemUserRole.Ctx(ctx).Insert(batch)
	}
	return err
}

func (s *sSysUser) FindByUsername(ctx context.Context, username string) (*sysentity.SysUser, error) {
	var u sysentity.SysUser
	err := dao.SystemUser.Ctx(ctx).
		Wheref("LOWER(username) = ?", username).
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
