package rbac

import "github.com/chibx/vuecom/backend/shared/models/db/users"

type Role struct {
	UserId      uint
	Name        string
	ParentID    *uint
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

func RoleFromBackend(backendRole *users.BackendRole, userId uint, excludedPerms ...string) *Role {
	role := &Role{
		Name:     backendRole.Name,
		ParentID: backendRole.ParentID,
		UserId:   userId,
	}

	_map := make(PermissionSet)
	if backendRole.AllowedPerms != nil {
		for _, v := range backendRole.AllowedPerms {
			perm := Permission(v)
			_map[perm] = struct{}{}
		}

		for _, v := range excludedPerms {
			delete(_map, Permission(v))
		}
	}

	role.permissions = _map

	return role
}
