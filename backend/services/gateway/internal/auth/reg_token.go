package auth

import (
	"fmt"
	"time"

	"github.com/chibx/vuecom/backend/services/gateway/internal/constants"
	"github.com/chibx/vuecom/backend/services/gateway/internal/global"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"
	"github.com/chibx/vuecom/backend/shared/errors/server"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type RegTokenField struct {
	RegToken string `json:"reg_token"`
	jwt.RegisteredClaims
}

func GenerateRegToken(api *types.Api, regToken *RegTokenField) (string, error) {
	regToken.RegisteredClaims = jwt.RegisteredClaims{
		Issuer:    api.AppName,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(constants.RegTokenDur)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, regToken)

	tokenString, err := token.SignedString(api.Config.SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateRegToken(api *types.Api, tokenString string, secretKey []byte) (RegTokenField, error) {
	logger := global.Logger
	var regToken = RegTokenField{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		logger.Error("Failed to validate backend jwt token", zap.Error(err))
		return regToken, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claimsToRegTokenField(claims, &regToken)
		return regToken, nil
	}

	return regToken, server.NewServerErr(fiber.StatusBadRequest, "Invalid token!")
}

func claimsToRegTokenField(claims jwt.MapClaims, regTokenStruc *RegTokenField) {
	for k, v := range claims {
		switch k {
		case "reg_token":
			regToken, _ := v.(string)
			regTokenStruc.RegToken = regToken
		case "iss":
			issuer, _ := v.(string)
			regTokenStruc.Issuer = issuer
		case "exp":
			// expiryStr, _ := strconv.ParseInt(v.(int64), 10, 64)
			expiryInt, _ := v.(int64)
			expiry := time.Unix(expiryInt, 0)
			regTokenStruc.ExpiresAt = jwt.NewNumericDate(expiry)
		case "iat":
			// expiryStr, _ := strconv.ParseInt(v.(int64), 10, 64)
			issuedAtInt, _ := v.(int64)
			issuedAt := time.Unix(issuedAtInt, 0)
			regTokenStruc.ExpiresAt = jwt.NewNumericDate(issuedAt)
		}
	}
}
