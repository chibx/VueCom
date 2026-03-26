package rbac

import (
	"slices"

	"github.com/chibx/vuecom/backend/shared/models/db/users"
)

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

func MergePermissions(_default, additional, excluded []string) PermissionSet {
	permMap := make(PermissionSet)

	for _, v := range _default {
		permMap[v] = struct{}{}
	}

	for _, v := range additional {
		permMap[v] = struct{}{}
	}

	for _, v := range excluded {
		delete(permMap, v)
	}

	return permMap
}

// IsValid checks if permission exists in registry
func IsValid(p Permission) bool {
	return slices.Contains(AllPermissions, p)
}
