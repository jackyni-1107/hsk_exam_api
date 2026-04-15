package member

import (
	"context"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"

	"github.com/gogf/gf/v2/errors/gerror"
	"golang.org/x/crypto/bcrypt"
)

func (s *sMember) MemberList(ctx context.Context, page, size int, username string, status int) ([]sysentity.SysMember, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
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
	err = model.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *sMember) MemberCreate(ctx context.Context, username, password, nickname, email, mobile, creator string, status int) (int64, error) {
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

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}

	if status != consts.StatusNormal && status != consts.StatusDisabled {
		status = consts.StatusNormal
	}

	id, err := dao.SysMember.Ctx(ctx).InsertAndGetId(sysdo.SysMember{
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
	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sMember) MemberUpdate(ctx context.Context, id int64, password, nickname, email, mobile, updater string, status int) error {
	var before sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return err
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeUserNotFound)
	}
	data := sysdo.SysMember{
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
	_, err := dao.SysMember.Ctx(ctx).Where("id", id).Data(data).Update()
	if err != nil {
		return err
	}
	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, &before, &after)
	}
	return nil
}

func (s *sMember) MemberDelete(ctx context.Context, id int64, updater string) error {
	var before sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return err
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeUserNotFound)
	}
	_, err := dao.SysMember.Ctx(ctx).Where("id", id).Data(sysdo.SysMember{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return err
	}
	var after sysentity.SysMember
	if err := dao.SysMember.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SysMember.Table(), id, &before, &after)
	}
	return nil
}

func (s *sMember) MemberProfile(ctx context.Context, memberId int64) (*sysentity.SysMember, error) {
	var m sysentity.SysMember
	err := dao.SysMember.Ctx(ctx).Where("id", memberId).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&m)
	if err != nil {
		return nil, err
	}
	if m.Id == 0 {
		return nil, gerror.NewCode(consts.CodeUserNotFound)
	}
	return &m, nil
}

func (s *sMember) FindByUsername(ctx context.Context, username string) (*sysentity.SysMember, error) {
	var m sysentity.SysMember
	err := dao.SysMember.Ctx(ctx).
		Where("username", username).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&m)
	if err != nil {
		return nil, err
	}
	if m.Id == 0 {
		return nil, nil
	}
	return &m, nil
}
