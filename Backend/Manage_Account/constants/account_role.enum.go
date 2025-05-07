package constants

import "slices"

// enum type
type AccountRole string

// enum values
const (
	AccountRoleGeneralUser AccountRole = "GENERAL_USER"
	AccountRoleAdmin       AccountRole = "ADMIN"
	AccountRoleSuperAdmin  AccountRole = "SUPER_ADMIN"
)

var ValidatorAccountRole []AccountRole = []AccountRole{
	AccountRoleGeneralUser,
	AccountRoleAdmin,
	AccountRoleSuperAdmin,
}

func IsAccountRoleExist(role AccountRole) bool {
	return slices.Contains(ValidatorAccountRole, role) == true
}

func GetAllAccountRole() []string {
	var role []string
	for i := range ValidatorAccountRole {
		role = append(role, string(ValidatorAccountRole[i]))
	}
	return role
}
