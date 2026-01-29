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

// Base Backend Panel User
type CreateBackendUserRequest struct {
	FullName    string  `json:"full_name" form:"full_name" validate:"required,min=5" name:"Full Name"`
	UserName    *string `json:"user_name" form:"user_name" validate:"required,min=3" name:"Username"`
	Email       string  `json:"email" form:"email" validate:"required,email" name:"Email Address"`
	PhoneNumber *string `json:"phone_number" form:"phone_number" validate:"min=10,max=15"`
	Country     *string `json:"country" form:"country" validate:"required_if=Role owner"`
	Password    string  `json:"password" form:"password" validate:"required,min=8,max=25"`
	// Role            string  `form:"role" validate:"required"`
	// "image" field for user logo (optional)
}

func (req *CreateBackendUserRequest) Validate() error {
	return utils.Validator().Struct(req)
}

func (req *CreateBackendUserRequest) ToDBBackendUser(ctx context.Context, api *types.Api, c *fiber.Ctx) (*userModels.BackendUser, error) {
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
		hashedUsername = *req.UserName // There is no need to encrypt the username (App specific)
	}

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

	// TODO: Implement image upload

	user := &userModels.BackendUser{
		FullName:     hashedFullname,
		UserName:     &hashedUsername,
		Email:        hashedEmail,
		PhoneNumber:  &hashedPhoneNumber,
		Image:        nil,
		PasswordHash: passwordHash,
		// TODO: I need to have a way to lookup a secure token (sent to the user through email) in the request url
		// c.Query("login_token"), then delete the token from the database,
		// instead of this
		// Role:            req.Role,
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
