package role

import rolesvc "exam/internal/service/role"

func init() {
	rolesvc.RegisterRole(new(sRole))
}
