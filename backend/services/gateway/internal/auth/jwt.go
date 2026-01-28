// Have chosen to go with Redis for session management. (for now)
package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"go.uber.org/zap"

	"github.com/golang-jwt/jwt/v5"
)

type JWTField struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// expirationTime is the time in seconds until the token expires.
func GenerateBackendAccessToken(api *types.Api, userId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTField{
		UserID: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    api.AppName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.BackendAccessTkDur)),
		},
	})

	tokenString, err := token.SignedString(api.Config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateCustomerAccessToken(api *types.Api, customerId int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &JWTField{
		UserID: customerId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    api.AppName,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.CustomerAccessTkDur)),
		},
	})

	tokenString, err := token.SignedString(api.Config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateBackendAccessToken(api *types.Api, tokenString string, secretKey []byte) (int, error) {
	logger := api.Deps.Logger
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		logger.Error("Failed to validate backend jwt token", zap.Error(err))
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var jwtField = JWTField{}
		claimsToJWTField(claims, &jwtField)
		userId := jwtField.UserID
		return userId, nil
	}

	return 0, fmt.Errorf("invalid token")
}

func ValidateCustomerAccessToken(api *types.Api, tokenString string, secretKey []byte) (int, error) {
	logger := api.Deps.Logger
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		logger.Error("Failed to validate customer jwt token", zap.Error(err))
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var jwtField = JWTField{}
		claimsToJWTField(claims, &jwtField)
		isValid := validateJWTMeta(api, &jwtField)
		if !isValid {
			return 0, errors.New("invalid session token")
		}
		userId := jwtField.UserID
		return userId, nil
	}

	return 0, fmt.Errorf("invalid token")
}

func claimsToJWTField(claims jwt.MapClaims, jwtField *JWTField) {
	for k, v := range claims {
		switch k {
		case "user_id":
			userId, _ := v.(int)
			jwtField.UserID = userId
		case "iss":
			issuer, _ := v.(string)
			jwtField.Issuer = issuer
		case "exp":
			// expiryStr, _ := strconv.ParseInt(v.(int64), 10, 64)
			expiryInt, _ := v.(int64)
			expiry := time.Unix(expiryInt, 0)
			jwtField.ExpiresAt = jwt.NewNumericDate(expiry)
		case "iat":
			// expiryStr, _ := strconv.ParseInt(v.(int64), 10, 64)
			issuedAtInt, _ := v.(int64)
			issuedAt := time.Unix(issuedAtInt, 0)
			jwtField.ExpiresAt = jwt.NewNumericDate(issuedAt)
		}
	}
}

func validateJWTMeta(api *types.Api, jwtStruct *JWTField) bool {
	if jwtStruct.Issuer != api.AppName {
		return false
	}

	if time.Now().After(jwtStruct.ExpiresAt.Time) {
		return false
	}

	return true
}
