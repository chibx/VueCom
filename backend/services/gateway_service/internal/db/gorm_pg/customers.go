package gorm_pg

import (
	"context"

	userModels "vuecom/shared/models/db/users"

	"gorm.io/gorm"
)

type customerRepository struct {
	db *gorm.DB
}

func (c *customerRepository) CreateUser(user *userModels.Customer, ctx context.Context) error {
	return c.db.WithContext(ctx).Create(user).Error
}

func (c *customerRepository) GetUserById(id int, ctx context.Context) (*userModels.Customer, error) {
	customer := &userModels.Customer{}
	err := c.db.WithContext(ctx).First(customer, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (c *customerRepository) GetSessionByToken(token string, ctx context.Context) (*userModels.CustomerSession, error) {
	sessionData := &userModels.CustomerSession{}

	err := c.db.WithContext(ctx).First(sessionData, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (c *customerRepository) GetSessions(customerId int, ctx context.Context) ([]userModels.CustomerSession, error) {
	var sessions []userModels.CustomerSession

	err := c.db.WithContext(ctx).Find(&sessions, "customer_id = ?", customerId).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

// CreateSession(session *dbModels.BackendSession, ctx context.Context) error
// 	DeleteSession(token string, ctx context.Context) error

func (c *customerRepository) CreateSession(session *userModels.CustomerSession, ctx context.Context) error {
	err := c.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *customerRepository) DeleteSession(session *userModels.CustomerSession, ctx context.Context) error {
	err := c.db.WithContext(ctx).Where("customer_id = ? AND token = ?", session.CustomerID, session.Token).Delete(session).Error
	if err != nil {
		return err
	}
	return nil
}
