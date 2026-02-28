package gorm_pg

import (
	"context"
	"time"

	// "strings"

	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/dto"
	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type backendUserRepository struct {
	db *gorm.DB
}

func (br *backendUserRepository) CreateRegToken(ctx context.Context, token string, supervisor uint, code string) error {
	var now = time.Now()
	var tokenStruc = &userModels.SignupToken{
		Token:      token,
		Code:       code,
		Supervisor: supervisor,
		CreatedAt:  now,
		ExpiryAt:   now.Add(constants.BackendRegTkDur),
	}

	return br.db.WithContext(ctx).Create(tokenStruc).Error
}

func (br *backendUserRepository) GetRegToken(ctx context.Context, token string) (*userModels.SignupToken, error) {
	var tokenStruc = &userModels.SignupToken{}
	err := br.db.WithContext(ctx).Where("token = ?", token).Preload("Super").First(tokenStruc).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, serverErrors.ErrDBRecordNotFound
		}
		return nil, err
	}

	return tokenStruc, nil
}

func (br *backendUserRepository) CreateUser(ctx context.Context, user *userModels.BackendUser) error {
	return br.db.WithContext(ctx).Create(user).Error
}

func (br *backendUserRepository) GetAdmin(ctx context.Context) (*userModels.BackendUser, error) {
	admin := &userModels.BackendUser{}

	err := br.db.Select("role").Where("role = ?", constants.OWNER).First(admin).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, serverErrors.ErrDBRecordNotFound
		}
		return nil, err
	}

	return admin, nil
}

func (br *backendUserRepository) HasAdmin(ctx context.Context) (bool, error) {
	var count int64

	err := br.db.Model(&userModels.BackendUser{}).WithContext(ctx).Where("role_id = ?", constants.OWNER_ID).Count(&count).Error

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
		if err == gorm.ErrRecordNotFound {
			return nil, serverErrors.ErrDBRecordNotFound
		}
		return nil, err
	}

	return user, nil
}

func (br *backendUserRepository) GetUserById(ctx context.Context, id int) (*userModels.BackendUser, error) {
	backendUser := &userModels.BackendUser{}
	err := br.db.WithContext(ctx).First(backendUser, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, serverErrors.ErrDBRecordNotFound
		}
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
		if err == gorm.ErrRecordNotFound {
			return nil, serverErrors.ErrDBRecordNotFound
		}
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
	err := br.db.WithContext(ctx).Delete(&userModels.BackendSession{}, "refresh_token_hash = ? AND device_id = ?", session.RefreshTokenHash, session.DeviceId).Error
	if err != nil {
		return err
	}
	return nil
}

func (br *backendUserRepository) GetCountryIdByCode(ctx context.Context, code string) (uint, error) {
	var country userModels.Country
	err := br.db.WithContext(ctx).Omit(clause.Associations).Select("id").Where("code = ?", code).First(&country).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, serverErrors.ErrDBRecordNotFound
		}
		return 0, err
	}

	return country.ID, nil
}
