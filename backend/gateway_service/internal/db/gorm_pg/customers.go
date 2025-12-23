package gorm_pg

import (
	"context"

	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func (c *customerRepository) GetCustomerById(id int, ctx context.Context) (*dbModels.Customer, error) {
	return nil, types.ErrDbUnimplemented
}
