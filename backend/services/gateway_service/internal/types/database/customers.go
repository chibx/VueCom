package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type CustomerRepository interface {
	CreateUser(user *dbModels.Customer, ctx context.Context) error
	GetUserById(id int, ctx context.Context) (*dbModels.Customer, error)
	CreateSession(session *dbModels.CustomerSession, ctx context.Context) error
	GetSessionByToken(token string, ctx context.Context) (*dbModels.CustomerSession, error)
	GetSessions(userId int, ctx context.Context) ([]dbModels.CustomerSession, error)
	DeleteSession(session *dbModels.CustomerSession, ctx context.Context) error
}
