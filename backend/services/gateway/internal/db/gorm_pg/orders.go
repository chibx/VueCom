package gorm_pg

import (
	"context"

	orderModels "github.com/chibx/vuecom/backend/shared/models/db/orders"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func (o *orderRepository) GetOrderById(ctx context.Context, id int) (*orderModels.Order, error) {
	order := &orderModels.Order{}

	err := o.db.WithContext(ctx).Where("id = ?", id).First(order).Error
	if err != nil {
		return nil, err
	}

	return order, nil
}

func (o *orderRepository) CreateOrder(ctx context.Context, order *orderModels.Order) error {
	return o.db.WithContext(ctx).Create(order).Error
}
