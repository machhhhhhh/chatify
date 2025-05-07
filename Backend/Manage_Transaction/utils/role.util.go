package utils

import (
	"chatify/constants"
	"slices"
)

func IsHavePermission(role constants.AccountRole, all_permission []constants.AccountRole) bool {
	return slices.Contains(all_permission, role) == true
}

func IsAdmin(role constants.AccountRole) bool {
	return role == constants.AccountRoleAdmin || role == constants.AccountRoleSuperAdmin
}

func IsSuperAdmin(role constants.AccountRole) bool {
	return role == constants.AccountRoleSuperAdmin
}
