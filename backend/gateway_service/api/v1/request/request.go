package request

import (
	"path"
	"vuecom/gateway/internal/v1/types"
)

type OnlyID struct {
	ID int `json:"id" params:"id"`
}

var IMAGE_FORMATS = []string{"image/jpeg", "image/png", "image/jpg"}
var GetBackendFolder = func(api *types.Api) string {
	// return api.AppName + "/backend_users"
	return path.Join(api.AppName, "backend_users")
}

var GetCustomerFolder = func(api *types.Api) string {
	return path.Join(api.AppName, "customers")
}

var GetProductFolder = func(api *types.Api) string {
	return path.Join(api.AppName, "products")
}
