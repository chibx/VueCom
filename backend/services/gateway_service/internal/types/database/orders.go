package database

import (
	"context"
	orderModels "vuecom/shared/models/db/orders"
)

type OrderRepository interface {
	CreateOrder(order *orderModels.Order, ctx context.Context) error
	GetOrderById(id int, ctx context.Context) (*orderModels.Order, error)
}
