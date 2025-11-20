package cache

import v1 "vuecom/gateway/api/v1"

func GetProduct(api v1.Api, id int) {
	_ = api.Redis
	_ = api.DB

}
