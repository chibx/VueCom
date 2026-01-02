package gorm_pg

import (
	"context"
	"errors"
	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type backendUserRepository struct {
	db *gorm.DB
}

func (br *backendUserRepository) CreateUser(user *dbModels.BackendUser, ctx context.Context) error {
	return br.db.WithContext(ctx).Create(user).Error
}

func (br *backendUserRepository) GetAdmin(ctx context.Context) (*dbModels.BackendUser, error) {
	admin := &dbModels.BackendUser{}

	err := br.db.Select("role").Where("role = 'owner'").First(admin).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrDbNil
		}

		return nil, err
	}

	return admin, nil
}

func (br *backendUserRepository) GetUserById(id int, ctx context.Context) (*dbModels.BackendUser, error) {
	backendUser := &dbModels.BackendUser{}
	err := br.db.WithContext(ctx).First(backendUser, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrDbNil
		}

		return nil, err
	}

	return backendUser, nil
}

func (br *backendUserRepository) GetUserByApiKey(apiKey string, ctx context.Context) (*dbModels.BackendUser, error) {
	return nil, types.ErrDbUnimplemented
}

func (br *backendUserRepository) GetSessionByToken(token string, ctx context.Context) (*dbModels.BackendSession, error) {
	sessionData := &dbModels.BackendSession{}

	err := br.db.WithContext(ctx).First(sessionData, "token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, types.ErrDbNil
		}

		return nil, err
	}

	return sessionData, nil
}

func (br *backendUserRepository) GetSessions(userId int, ctx context.Context) ([]dbModels.BackendSession, error) {
	var sessions []dbModels.BackendSession

	err := br.db.WithContext(ctx).Find(&sessions, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (br *backendUserRepository) CreateSession(session *dbModels.BackendSession, ctx context.Context) error {
	err := br.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (br *backendUserRepository) DeleteSession(session *dbModels.BackendSession, ctx context.Context) error {
	err := br.db.WithContext(ctx).Where("user_id = ? AND token = ?", session.UserId, session.Token).Delete(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (br *backendUserRepository) GetCountryIdByCode(code string, ctx context.Context) (uint, error) {
	var country dbModels.Country
	err := br.db.WithContext(ctx).Omit(clause.Associations).Select("id").Where("code = ?", code).First(&country).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, types.ErrDbNil
		}
		return 0, err
	}
	return country.ID, nil
}
