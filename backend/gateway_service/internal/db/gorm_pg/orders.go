package gorm_pg

import (
	"context"
	"errors"

	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func (o *orderRepository) GetOrderById(id int, ctx context.Context) (*dbModels.Order, error) {
	order := &dbModels.Order{}

	err := o.db.WithContext(ctx).Where("id = ?", id).First(order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrDbNil
		}
		return nil, err
	}

	return order, nil
}

func (o *orderRepository) CreateOrder(order *dbModels.Order, ctx context.Context) error {
	return o.db.WithContext(ctx).Create(order).Error
}
