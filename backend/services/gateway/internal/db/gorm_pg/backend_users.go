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

func (br *backendUserRepository) GetUserByNameForLogin(ctx context.Context, username string) (*dto.UserForLogin, error) {
	// selectedValue := "*"
	// backendUser := &userModels.BackendUser{}

	// if fields != nil {
	// 	selectedValue = strings.Join(fields, ",")
	// }
	user := &dto.UserForLogin{}

	// userModels.BackendUser
	err := br.db.Model("backend_users").Where("user_name = ?", username).First(user).Error

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

func (br *backendUserRepository) GetSessionByToken(ctx context.Context, token string) (*userModels.BackendSession, error) {
	sessionData := &userModels.BackendSession{}

	err := br.db.WithContext(ctx).First(sessionData, "token = ?", token).Error
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

func (br *backendUserRepository) DeleteSession(session *userModels.BackendSession, ctx context.Context) error {
	err := br.db.WithContext(ctx).Where("user_id = ? AND token = ?", session.UserId, session.Token).Delete(session).Error
	if err != nil {
		return err
	}
	return nil
}

func (br *backendUserRepository) GetCountryIdByCode(code string, ctx context.Context) (uint, error) {
	var country userModels.Country
	err := br.db.WithContext(ctx).Omit(clause.Associations).Select("id").Where("code = ?", code).First(&country).Error
	if err != nil {
		return 0, err
	}
	return country.ID, nil
}
