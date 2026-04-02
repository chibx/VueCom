package rbac

type Permission = string
type PermissionSet map[Permission]struct{}

type Role struct {
	UserId      uint
	Name        string
	ParentID    *uint
	permissions PermissionSet
}

type UserRoleFromDB struct {
	RoleID          uint     `gorm:""`
	ExcludedPerms   []string `gorm:"column:excluded_permissions;type:text[]"`
	AdditionalPerms []string `gorm:"column:additional_permissions;type:text[]"`
}

func (role *Role) Has(perms ...Permission) bool {
	if len(perms) == 0 {
		return false
	}

	if len(perms) > 1 {
		hasAll := true
		for _, p := range perms {
			if _, ok := role.permissions[p]; !ok {
				hasAll = false
				break
			}
		}

		return hasAll
	}

	_, ok := role.permissions[perms[0]]

	return ok
}

func (set PermissionSet) Has(perms ...string) bool {
	if _, ok := set["*"]; ok {
		return true
	}

	if len(perms) == 0 {
		return false
	}

	if len(perms) > 1 {
		hasAll := true
		for _, v := range perms {
			if _, ok := set[v]; !ok {
				hasAll = false
				break
			}
		}

		return hasAll
	}

	_, ok := set[perms[0]]

	return ok
}
