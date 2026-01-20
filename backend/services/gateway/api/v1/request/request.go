package request

import (
	"path"

	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
)

const CUSTOMER_KEY = "customer"
const CUSTOMER_TOKEN = "c_token"
const BACKEND_USER_KEY = "backend_user"
const BACKEND_TOKEN = "b_token"

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
