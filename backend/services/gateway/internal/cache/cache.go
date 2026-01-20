package cache

import (
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
)

func GetProduct(api types.Api, id int) {
	_ = api.Deps.Redis
	_ = api.Deps.DB

}
