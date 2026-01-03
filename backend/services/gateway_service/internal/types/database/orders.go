package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type OrderRepository interface {
	CreateOrder(order *dbModels.Order, ctx context.Context) error
	GetOrderById(id int, ctx context.Context) (*dbModels.Order, error)
}
