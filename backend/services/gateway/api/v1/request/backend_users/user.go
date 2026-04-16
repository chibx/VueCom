package backendusers

import (
	"context"
	"errors"
	"fmt"

	userModels "github.com/chibx/vuecom/backend/shared/models/db/users"

	"github.com/chibx/vuecom/backend/services/gateway/internal/auth"
	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"

	serverErrors "github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type CreateOwnerRequest struct {
	FullName    string  `json:"full_name" form:"full_name" validate:"required,min=5" name:"Full Name"`
	UserName    string  `json:"user_name" form:"user_name" validate:"required,min=3" name:"Username"`
	Email       string  `json:"email" form:"email" validate:"required,email" name:"Email Address"`
	PhoneNumber *string `json:"phone_number" form:"phone_number" validate:"min=10,max=15"`
	Country     *string `json:"country" form:"country" validate:"required"`
	Password    string  `json:"password" form:"password" validate:"required,min=8,max=25"`
}

// Base Backend Panel User
type CreateBackendUserRequest struct {
	FullName    string  `json:"full_name" form:"full_name" validate:"required,min=5" name:"Full Name"`
	UserName    string  `json:"user_name" form:"user_name" validate:"required,min=3" name:"Username"`
	Code        string  `json:"code" form:"code" validate:"required" name:"Signup Code"`
	Email       string  `json:"email" form:"email" validate:"required,email" name:"Email Address"`
	PhoneNumber *string `json:"phone_number" form:"phone_number" validate:"min=10,max=15"`
	Country     *string `json:"country" form:"country" validate:"required"`
	Password    string  `json:"password" form:"password" validate:"required,min=6"`
}

type CreateTokenRequest struct {
	Supervisor uint   `json:"supervisor" validate:"required,gt=0"`
	Email      string `json:"email" validate:"required,email"`
}

type UserLogin struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func (req *CreateOwnerRequest) Validate() error {
	return utils.Validator().Struct(req)
}

func (req *CreateOwnerRequest) ToDBBackendUser(ctx context.Context, api *types.Api, c *fiber.Ctx) (*userModels.BackendUser, error) {
	db := api.Deps.DB
	logger := global.Logger

	passwordHash, err := auth.GenerateHashFromString(req.Password, auth.DefaultHashParams)
	if err != nil {
		logger.Error("Failed to hash password for new backend user", zap.Error(err))
		return nil, err
	}

	encryptedFullname, err := auth.Encrypt(req.FullName, api.Config.SecretKey)
	if err != nil {
		logger.Error("Failed to encrypt fullname for new backend user", zap.Error(err))
		return nil, err
	}
	var username = req.UserName
	// if req.UserName != nil {
	// 	username = *req.UserName // There is no need to encrypt the username (App specific)
	// }

	encryptedEmail, err := auth.Encrypt(req.Email, api.Config.SecretKey)
	if err != nil {
		logger.Error("Failed to encrypt email for new backend user", zap.Error(err))
		return nil, err
	}

	var encryptedPhoneNumber string
	if req.PhoneNumber != nil {
		encryptedPhoneNumber, err = auth.Encrypt(*req.PhoneNumber, api.Config.SecretKey)
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
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return nil, serverErrors.NewServerErr(fiber.StatusBadRequest, fmt.Sprintf("Record for Country %s does not exist", *req.Country))
			}
			return nil, err
		}
	}

	// TODO: Implement image upload
	var roleId = uint(constants.OWNER_ID)
	user := &userModels.BackendUser{
		FullName:     encryptedFullname,
		UserName:     username,
		Email:        encryptedEmail,
		PhoneNumber:  &encryptedPhoneNumber,
		Image:        nil,
		PasswordHash: passwordHash,
		// TODO: I need to have a way to lookup a secure token (sent to the user through email) in the request url
		// c.Query("login_token"), then delete the token from the database,
		// instead of this
		RoleID:          roleId,
		IsEmailVerified: true,
	}

	if req.PhoneNumber != nil {
		user.PhoneNumber = &encryptedPhoneNumber
	}
	if countryId != 0 {
		user.CountryId = &countryId
	}

	return user, nil
}

func (req *CreateBackendUserRequest) ToDBBackendUser(ctx context.Context, api *types.Api, c *fiber.Ctx) (*userModels.BackendUser, error) {
	db := api.Deps.DB
	logger := global.Logger

	passwordHash, err := auth.GenerateHashFromString(req.Password, auth.DefaultHashParams)
	if err != nil {
		logger.Error("Failed to hash password for new backend user", zap.Error(err))
		return nil, err
	}

	encryptedFullname, err := auth.Encrypt(req.FullName, api.Config.SecretKey)
	if err != nil {
		logger.Error("Failed to encrypt fullname for new backend user", zap.Error(err))
		return nil, err
	}
	var username = req.UserName
	// if req.UserName != nil {
	// 	username = *req.UserName // There is no need to encrypt the username (App specific)
	// }

	encryptedEmail, err := auth.Encrypt(req.Email, api.Config.SecretKey)
	if err != nil {
		logger.Error("Failed to encrypt email for new backend user", zap.Error(err))
		return nil, err
	}

	var encryptedPhoneNumber string
	if req.PhoneNumber != nil {
		encryptedPhoneNumber, err = auth.Encrypt(*req.PhoneNumber, api.Config.SecretKey)
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
			if errors.Is(err, serverErrors.ErrDBRecordNotFound) {
				return nil, serverErrors.NewServerErr(fiber.StatusBadRequest, fmt.Sprintf("Record for Country %s does not exist", *req.Country))
			}
			return nil, err
		}
	}

	// TODO: Implement image upload

	user := &userModels.BackendUser{
		FullName:        encryptedFullname,
		UserName:        username,
		Email:           encryptedEmail,
		PhoneNumber:     &encryptedPhoneNumber,
		Image:           nil,
		PasswordHash:    passwordHash,
		IsEmailVerified: true,

		// TODO: I need to have a way to lookup a secure token (sent to the user through email) in the request url
		// c.Query("login_token"), then delete the token from the database,
		// instead of this
		// Role: constants.OWNER,

	}

	// if req.UserName != nil {
	// 	user.UserName = &username
	// }
	if req.PhoneNumber != nil {
		user.PhoneNumber = &encryptedPhoneNumber
	}
	if countryId != 0 {
		user.CountryId = &countryId
	}

	return user, nil
}
