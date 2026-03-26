package cache

import (
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"github.com/chibx/vuecom/backend/shared/rbac"
	lru "github.com/hashicorp/golang-lru/v2"
)

// user_id -> permission[]
var RoleCache *lru.Cache[int, []string]
var UserPermCache *lru.Cache[int, rbac.PermissionSet]

func InitInMemCache() {
	var err error
	RoleCache, err = lru.New[int, []string](1000)
	if err != nil {
		utils.Logger().Fatal("Couldn't initialize role and permission in-memory cache")
	}
	UserPermCache, err = lru.New[int, rbac.PermissionSet](1000)
	if err != nil {
		utils.Logger().Fatal("Couldn't initialize role and permission in-memory cache")
	}
}
