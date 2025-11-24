package cache

import (
	"vuecom/gateway/internal/v1/types"
)

func GetProduct(api types.Api, id int) {
	_ = api.Deps.Redis
	_ = api.Deps.DB

}
