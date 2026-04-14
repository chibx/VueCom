package keys

import "strconv"

const APP_DATA_KEY = "app_data"
const CUSTOMER_KEY_PREFIX = "cust:"
const PRODUCT_KEY_PREFIX = "product:"
const BACKEND_USER_KEY_PREFIX = "b_user:"

func CustomerKey(userId uint32) string {
	return CUSTOMER_KEY_PREFIX + strconv.FormatUint(uint64(userId), 10)
}

func ProductKey(productId uint32) string {
	return PRODUCT_KEY_PREFIX + strconv.FormatUint(uint64(productId), 10)
}

func BackendUserKey(userId uint32) string {
	return BACKEND_USER_KEY_PREFIX + strconv.FormatUint(uint64(userId), 10)
}
