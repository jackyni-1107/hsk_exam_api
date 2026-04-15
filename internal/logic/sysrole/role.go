package sysrole

import (
	"context"

	"exam/internal/auditutil"
	"exam/internal/consts"
	"exam/internal/dao"
	sysdo "exam/internal/model/do/sys"
	sysentity "exam/internal/model/entity/sys"
	"exam/internal/utility"

	"github.com/gogf/gf/v2/errors/gerror"
)

func (s *sSysRole) RoleList(ctx context.Context, page, size int, name string, status int) ([]sysentity.SysRole, int, error) {
	if page <= 0 {
		page = 1
	}
	if size <= 0 {
		size = 10
	}
	model := dao.SystemRole.Ctx(ctx).Where("delete_flag", consts.DeleteFlagNotDeleted)
	if name != "" {
		model = model.WhereLike("name", "%"+name+"%")
	}
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		model = model.Where("status", status)
	}
	total, err := model.Count()
	if err != nil {
		return nil, 0, err
	}
	var list []sysentity.SysRole
	err = model.Page(page, size).OrderDesc("id").Scan(&list)
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (s *sSysRole) RoleMenuIds(ctx context.Context, roleId int64) ([]int64, error) {
	var rows []sysentity.SysRoleMenu
	err := dao.SystemRoleMenu.Ctx(ctx).
		Where("role_id", roleId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Scan(&rows)
	if err != nil {
		return nil, err
	}
	ids := make([]int64, 0, len(rows))
	for _, r := range rows {
		ids = append(ids, r.MenuId)
	}
	return ids, nil
}

func (s *sSysRole) RoleCreate(ctx context.Context, name, code, remark, creator string, status, sort, typ int) (int64, error) {
	cnt, err := dao.SystemRole.Ctx(ctx).
		Where("code", code).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Count()
	if err != nil {
		return 0, err
	}
	if cnt > 0 {
		return 0, gerror.NewCode(consts.CodeRoleCodeExists)
	}

	if status != consts.StatusNormal && status != consts.StatusDisabled {
		status = consts.StatusNormal
	}

	id, err := dao.SystemRole.Ctx(ctx).InsertAndGetId(sysdo.SysRole{
		Name:    name,
		Code:    code,
		Remark:  remark,
		Sort:    sort,
		Status:  status,
		Type:    typ,
		Creator: creator,
		Updater: creator,
	})
	if err != nil {
		return 0, err
	}
	var after sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Scan(&after); err == nil && after.Id > 0 {
		auditutil.RecordEntityDiff(ctx, dao.SystemRole.Table(), id, nil, &after)
	}
	return id, nil
}

func (s *sSysRole) RoleUpdate(ctx context.Context, id int64, name, code, remark, updater string, status, sort, typ int) error {
	var before sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return err
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeRoleNotFound)
	}
	data := sysdo.SysRole{
		Name:    name,
		Code:    code,
		Remark:  remark,
		Sort:    sort,
		Type:    typ,
		Updater: updater,
	}
	if status == consts.StatusNormal || status == consts.StatusDisabled {
		data.Status = status
	}
	_, err := dao.SystemRole.Ctx(ctx).Where("id", id).Data(data).Update()
	if err != nil {
		return err
	}
	var after sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemRole.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysRole) RoleDelete(ctx context.Context, id int64, updater string) error {
	var before sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Where("delete_flag", consts.DeleteFlagNotDeleted).Scan(&before); err != nil {
		return err
	}
	if before.Id == 0 {
		return gerror.NewCode(consts.CodeRoleNotFound)
	}
	_, err := dao.SystemRole.Ctx(ctx).Where("id", id).Data(sysdo.SysRole{
		DeleteFlag: consts.DeleteFlagDeleted,
		Updater:    updater,
	}).Update()
	if err != nil {
		return err
	}
	var after sysentity.SysRole
	if err := dao.SystemRole.Ctx(ctx).Where("id", id).Scan(&after); err == nil {
		auditutil.RecordEntityDiff(ctx, dao.SystemRole.Table(), id, &before, &after)
	}
	return nil
}

func (s *sSysRole) RoleMenuAssign(ctx context.Context, roleId int64, menuIds []int64, creator string) error {
	beforeIDs, err := s.RoleMenuIds(ctx, roleId)
	if err != nil {
		return err
	}
	beforeStr := utility.JoinSortedInt64IDs(beforeIDs)
	_, err = dao.SystemRoleMenu.Ctx(ctx).
		Where("role_id", roleId).
		Where("delete_flag", consts.DeleteFlagNotDeleted).
		Data(sysdo.SysRoleMenu{
			DeleteFlag: consts.DeleteFlagDeleted,
			Updater:    creator,
		}).Update()
	if err != nil {
		return err
	}

	if len(menuIds) > 0 {
		batch := make([]sysdo.SysRoleMenu, 0, len(menuIds))
		for _, mid := range menuIds {
			batch = append(batch, sysdo.SysRoleMenu{
				RoleId:  roleId,
				MenuId:  mid,
				Creator: creator,
				Updater: creator,
			})
		}
		_, err = dao.SystemRoleMenu.Ctx(ctx).Insert(batch)
	}
	if err != nil {
		return err
	}
	afterStr := utility.JoinSortedInt64IDs(menuIds)
	auditutil.RecordMapDiff(ctx, dao.SystemRoleMenu.Table(), roleId,
		map[string]interface{}{"menu_ids": beforeStr},
		map[string]interface{}{"menu_ids": afterStr})
	return nil
}
