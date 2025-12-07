package keys

import "strconv"

const APP_DATA_KEY = "app_data"
const USER_KEY_PREFIX = "user:"
const PRODUCT_KEY_PREFIX = "product:"
const BACKEND_USER_KEY_PREFIX = "backend_user:"

func UserKey(userId int) string {
	return USER_KEY_PREFIX + strconv.Itoa(userId)
}

func ProductKey(productId int) string {
	return PRODUCT_KEY_PREFIX + strconv.Itoa(productId)
}

func BackendUserKey(userId int) string {
	return BACKEND_USER_KEY_PREFIX + strconv.Itoa(userId)
}
