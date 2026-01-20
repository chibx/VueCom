package database

import (
	"context"
	userModels "vuecom/shared/models/db/users"
)

type CustomerRepository interface {
	CreateUser(user *userModels.Customer, ctx context.Context) error
	GetUserById(id int, ctx context.Context) (*userModels.Customer, error)
	CreateSession(session *userModels.CustomerSession, ctx context.Context) error
	GetSessionByToken(token string, ctx context.Context) (*userModels.CustomerSession, error)
	GetSessions(userId int, ctx context.Context) ([]userModels.CustomerSession, error)
	DeleteSession(session *userModels.CustomerSession, ctx context.Context) error
}
