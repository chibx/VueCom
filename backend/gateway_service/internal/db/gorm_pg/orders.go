package gorm_pg

import (
	"context"

	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type orderRepository struct {
	db *gorm.DB
}

func (o *orderRepository) GetOrderById(id int, ctx context.Context) (*dbModels.Order, error) {
	return nil, types.ErrDbUnimplemented
}
