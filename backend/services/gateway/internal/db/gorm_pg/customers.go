package gorm_pg

import (
	"context"

	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func (c *customerRepository) CreateUser(ctx context.Context, user *userModels.Customer) error {
	return c.db.WithContext(ctx).Create(user).Error
}

func (c *customerRepository) GetUserById(ctx context.Context, id int) (*userModels.Customer, error) {
	customer := &userModels.Customer{}
	err := c.db.WithContext(ctx).First(customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c *customerRepository) GetSessionByTokenId(ctx context.Context, tokenId string) (*userModels.CustomerSession, error) {
	sessionData := &userModels.CustomerSession{}

	err := c.db.WithContext(ctx).First(sessionData, "id = ?", tokenId).Error
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (c *customerRepository) GetSessions(ctx context.Context, customerId int) ([]userModels.CustomerSession, error) {
	var sessions []userModels.CustomerSession

	err := c.db.WithContext(ctx).Find(&sessions, "customer_id = ?", customerId).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// CreateSession(session *dbModels.BackendSession, ctx context.Context) error
// 	DeleteSession(token string, ctx context.Context) error

func (c *customerRepository) CreateSession(ctx context.Context, session *userModels.CustomerSession) error {
	err := c.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepository) DeleteSession(ctx context.Context, session *userModels.CustomerSession) error {
	err := c.db.WithContext(ctx).Delete(session).Error
	if err != nil {
		return err
	}
	return nil
}
