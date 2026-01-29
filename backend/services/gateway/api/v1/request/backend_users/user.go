package backendusers

import (
	"context"
	"errors"
	"fmt"

	"github.com/chibx/vuecom/backend/shared/errors/server"
	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// JWT Format sent back to the client dashboard
type BackendJWTPayload struct {
	UserId int    `json:"user_id"`
	Role   string `json:"role"`
}

// type sharedUserProps struct {
// 	FullName        string  `json:"full_name" gorm:"not null;type:varchar(255);index" validate:""`
// 	UserName        *string `json:"user_name" gorm:"type:varchar(255);index" validate:""`
// 	Email           string  `json:"email" gorm:"unique;not null;type:varchar(255);index"`
// 	PhoneNumber     *string `json:"phone_number" gorm:"type:varchar(20)" validate:""`
// 	Image           *string `json:"image" validate:"url"`
// 	Country         uint    `json:"country" gorm:"index"`
// 	IsEmailVerified bool    `json:"email_verified" gorm:"default:FALSE;not null"`
// 	Password        *string `json:"password,omitempty" validate:"required"`
// }

// Base Backend Panel User
type CreateBackendUserRequest struct {
	FullName        string  `form:"full_name" validate:"required,min=5" name:"Full Name"`
	UserName        *string `form:"user_name" validate:"required,min=3" name:"Username"`
	Email           string  `form:"email" validate:"required,email" name:"Email Address"`
	PhoneNumber     *string `form:"phone_number" validate:"min=10,max=15"`
	Country         *string `form:"country" validate:"required_if=Role owner"`
	IsEmailVerified bool    `form:"email_verified"`
	Password        string  `form:"password" validate:"required,min=8,max=25"`
	Role            string  `form:"role" validate:"required"`
}

func (req *CreateBackendUserRequest) Validate() error {
	return utils.Validator().Struct(req)
}

func (req *CreateBackendUserRequest) ToDBBackendUser(api *types.Api, ctx context.Context) (*userModels.BackendUser, error) {
	db := api.Deps.DB
	logger := utils.Logger()

	passwordHash, err := auth.GenerateHashFromString(req.Password, auth.DefaultHashParams)
	if err != nil {
		logger.Error("Failed to hash password for new backend user", zap.Error(err))
		return nil, err
	}

	hashedFullname, err := auth.Encrypt(req.FullName, api.Config.SecretKey)
	if err != nil {
		logger.Error("Failed to encrypt fullname for new backend user", zap.Error(err))
		return nil, err
	}
	var hashedUsername string
	if req.UserName != nil {
		// hashedUsername, err = auth.Encrypt(*req.UserName, api.Config.DbEncKey)
		hashedUsername = *req.UserName // There is no need to encrypt the username (App specific)
	}
	// if err != nil {
	// 	return nil, err
	// }

	hashedEmail, err := auth.Encrypt(req.Email, api.Config.SecretKey)
	if err != nil {
		logger.Error("Failed to encrypt email for new backend user", zap.Error(err))
		return nil, err
	}

	var hashedPhoneNumber string
	if req.PhoneNumber != nil {
		hashedPhoneNumber, err = auth.Encrypt(*req.PhoneNumber, api.Config.SecretKey)
	}
	if err != nil {
		logger.Error("Failed to encrypt phone number for new backend user", zap.Error(err))
		return nil, err
	}

	var countryId uint
	if req.Country != nil {
		// err = api.Deps.DB.Model(&dbModels.Country{}).Where(dbModels.Country{Code: *req.Country}).Row().Scan(&countryId)
		countryId, err = db.BackendUsers().GetCountryIdByCode(ctx, *req.Country)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get country ID for `%s` for new backend user", *req.Country), zap.Error(err))
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, server.NewServerErr(fiber.StatusBadRequest, fmt.Sprintf("Record for Country %s does not exist", *req.Country))
			}
			return nil, err
		}
	}

	// TODO: Implement
	// TODO: Implement image upload

	user := &userModels.BackendUser{
		FullName:        hashedFullname,
		UserName:        &hashedUsername,
		Email:           hashedEmail,
		PhoneNumber:     &hashedPhoneNumber,
		Image:           nil,
		IsEmailVerified: req.IsEmailVerified,
		PasswordHash:    passwordHash,
		Role:            req.Role,
	}

	if req.UserName != nil {
		user.UserName = &hashedUsername
	}
	if req.PhoneNumber != nil {
		user.PhoneNumber = &hashedPhoneNumber
	}
	if countryId != 0 {
		user.CountryId = &countryId
	}

	return user, nil
}
