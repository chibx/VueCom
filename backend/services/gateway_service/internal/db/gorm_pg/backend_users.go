package gorm_pg

import (
	"context"
	userModels "vuecom/shared/models/db/users"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type backendUserRepository struct {
	db *gorm.DB
}

func (br *backendUserRepository) CreateUser(user *userModels.BackendUser, ctx context.Context) error {
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

func (br *backendUserRepository) GetUserById(id int, ctx context.Context) (*userModels.BackendUser, error) {
	backendUser := &userModels.BackendUser{}
	err := br.db.WithContext(ctx).First(backendUser, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return backendUser, nil
}

func (br *backendUserRepository) GetUserByApiKey(apiKey string, ctx context.Context) (*userModels.BackendUser, error) {
	return nil, errDbUnimplemented
}

func (br *backendUserRepository) GetSessionByToken(token string, ctx context.Context) (*userModels.BackendSession, error) {
	sessionData := &userModels.BackendSession{}

	err := br.db.WithContext(ctx).First(sessionData, "token = ?", token).Error
	if err != nil {
		return nil, err
	}

	return sessionData, nil
}

func (br *backendUserRepository) GetSessions(userId int, ctx context.Context) ([]userModels.BackendSession, error) {
	var sessions []userModels.BackendSession

	err := br.db.WithContext(ctx).Find(&sessions, "user_id = ?", userId).Error
	if err != nil {
		return nil, err
	}

	return sessions, nil
}

func (br *backendUserRepository) CreateSession(session *userModels.BackendSession, ctx context.Context) error {
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
