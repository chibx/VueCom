package database

import (
	"context"
	dbModels "vuecom/shared/models/db"
)

type BackendUserRepository interface {
	CreateUser(user *dbModels.BackendUser, ctx context.Context) error
	GetUserById(id int, ctx context.Context) (*dbModels.BackendUser, error)
	GetAdmin(ctx context.Context) (*dbModels.BackendUser, error)
	// The full api key is passed to the function
	//
	// Splitting should be done by the function
	GetUserByApiKey(apiKey string, ctx context.Context) (*dbModels.BackendUser, error)
	CreateSession(session *dbModels.BackendSession, ctx context.Context) error
	GetSessionByToken(token string, ctx context.Context) (*dbModels.BackendSession, error)
	GetSessions(userId int, ctx context.Context) ([]dbModels.BackendSession, error)
	DeleteSession(session *dbModels.BackendSession, ctx context.Context) error

	GetCountryIdByCode(code string, ctx context.Context) (uint, error)
}
