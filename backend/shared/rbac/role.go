package rbac

type Role struct {
	UserId      uint
	Name        string
	ParentID    *int
	permissions PermissionSet
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
