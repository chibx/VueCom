package gorm_pg

import (
	"context"
	// "strings"

	"github.com/chibx/vuecom/backend/services/gateway/internal/dto"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type backendUserRepository struct {
	db *gorm.DB
}

func (br *backendUserRepository) CreateUser(ctx context.Context, user *userModels.BackendUser) error {
	return br.db.WithContext(ctx).Create(user).Error
}

func (br *backendUserRepository) GetAdmin(ctx context.Context) (*userModels.BackendUser, error) {
	admin := &userModels.BackendUser{}

	err := br.db.Select("role").Where("role = 'owner'").First(admin).Error
	if err != nil {
		return nil, err
	}

	return admin, nil
}

func (br *backendUserRepository) HasAdmin(ctx context.Context) (bool, error) {
	var count int64

	err := br.db.Where("role = 'owner'").Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (br *backendUserRepository) GetUserByNameForLogin(ctx context.Context, username string) (*dto.UserForLogin, error) {
	var user *dto.UserForLogin

	// userModels.BackendUser
	err := br.db.Model(&userModels.BackendUser{}).Where("user_name = ?", username).First(user).Error

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (br *backendUserRepository) GetUserById(ctx context.Context, id int) (*userModels.BackendUser, error) {
	backendUser := &userModels.BackendUser{}
	err := br.db.WithContext(ctx).First(backendUser, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return backendUser, nil
}

func (br *backendUserRepository) GetUserByApiKey(ctx context.Context, apiKey string) (*userModels.BackendUser, error) {
	return nil, errDbUnimplemented
}

func (br *backendUserRepository) GetSessionByTokenId(ctx context.Context, tokenId string) (*userModels.BackendSession, error) {
	sessionData := &userModels.BackendSession{}

	err := br.db.WithContext(ctx).First(sessionData, "id = ?", tokenId).Error
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (br *backendUserRepository) GetSessions(ctx context.Context, userId int) ([]userModels.BackendSession, error) {
	var sessions []userModels.BackendSession

	err := br.db.WithContext(ctx).Find(&sessions, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (br *backendUserRepository) CreateSession(ctx context.Context, session *userModels.BackendSession) error {
	err := br.db.WithContext(ctx).Create(session).Error
	if err != nil {
		return err
	}

	return nil
}

func (br *backendUserRepository) DeleteSession(ctx context.Context, session *userModels.BackendSession) error {
	err := br.db.WithContext(ctx).Delete(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (br *backendUserRepository) GetCountryIdByCode(ctx context.Context, code string) (uint, error) {
	var country userModels.Country
	err := br.db.WithContext(ctx).Omit(clause.Associations).Select("id").Where("code = ?", code).First(&country).Error
	if err != nil {
		return 0, err
	}
	return country.ID, nil
}
