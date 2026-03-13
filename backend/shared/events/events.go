package events

type eventType string

// Catalog
const (
	PRODUCT_CREATION eventType = "product:creation"
	PRODUCT_UPDATE   eventType = "product:update"
	PRODUCT_DELETE   eventType = "product:delete"
	CATEGORY_CHANGE  eventType = "category:change"
	REVIEW_SUBMITTED eventType = "review:submitted"
)

// Orders
const (
	ORDER_CREATED   eventType = "order:created"
	ORDER_UPDATED   eventType = "order:updated"
	ORDER_CANCELLED eventType = "order:cancelled"
	ORDER_SHIPPED   eventType = "order:shipped"
	ORDER_DELIVERED eventType = "order:delivered"
	ORDER_REFUNDED  eventType = "order:refunded"
)

// Notifications
const (
	NOTIFICATION_SEND eventType = "notification:send"
)

// Inventory
const (
	INVENTORY_LOW_STOCK eventType = "inventory:low_stock"
	INVENTORY_UPDATED   eventType = "inventory:updated"
	BACKORDER_CREATED   eventType = "backorder:created"
)

// Analytics
const (
	PAGE_VIEW         eventType = "analytics:page_view"
	ADD_TO_CART       eventType = "analytics:add_to_cart"
	CHECKOUT_START    eventType = "analytics:checkout_start"
	PURCHASE_COMPLETE eventType = "analytics:purchase_complete"
)

// Email
const (
	EMAIL_SEND eventType = "email:send"
)

// Payment
const (
	PAYMENT_SUCCEEDED eventType = "payment:succeeded"
	PAYMENT_FAILED    eventType = "payment:failed"
	PAYMENT_REFUNDED  eventType = "payment:refunded"
)
