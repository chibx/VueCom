package cache

import (
	"vuecom/gateway/internal/types"
)

func GetProduct(api types.Api, id int) {
	_ = api.Deps.Redis
	_ = api.Deps.DB

}
