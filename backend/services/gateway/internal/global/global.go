package global

import (
	"github.com/chibx/vuecom/backend/shared/rbac"
	lru "github.com/hashicorp/golang-lru/v2"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func SetLogger(l *zap.Logger) {
	Logger = l
}

// user_id -> permission[]
var RoleCache *lru.Cache[int, []string]
var UserPermCache *lru.Cache[int, rbac.PermissionSet]

func InitInMemCache() {
	var err error
	RoleCache, err = lru.New[int, []string](1000)
	if err != nil {
		Logger.Fatal("Couldn't initialize role and permission in-memory cache")
	}
	UserPermCache, err = lru.New[int, rbac.PermissionSet](1000)
	if err != nil {
		Logger.Fatal("Couldn't initialize role and permission in-memory cache")
	}
}
