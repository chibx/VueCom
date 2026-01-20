package database

import (
	"context"
	userModels "vuecom/shared/models/db/users"
)

type BackendUserRepository interface {
	CreateUser(user *userModels.BackendUser, ctx context.Context) error
	GetUserById(id int, ctx context.Context) (*userModels.BackendUser, error)
	GetAdmin(ctx context.Context) (*userModels.BackendUser, error)
	// The full api key is passed to the function
	//
	// Splitting should be done by the function
	GetUserByApiKey(apiKey string, ctx context.Context) (*userModels.BackendUser, error)
	CreateSession(session *userModels.BackendSession, ctx context.Context) error
	GetSessionByToken(token string, ctx context.Context) (*userModels.BackendSession, error)
	GetSessions(userId int, ctx context.Context) ([]userModels.BackendSession, error)
	DeleteSession(session *userModels.BackendSession, ctx context.Context) error

	GetCountryIdByCode(code string, ctx context.Context) (uint, error)
}
