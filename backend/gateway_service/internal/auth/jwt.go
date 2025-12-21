// Have chosen to go with Redis for session management. (for now)
package auth

import (
	"fmt"
	"vuecom/gateway/internal/v1/types"

	"github.com/golang-jwt/jwt/v5"
)

type JWTField struct {
	UserName string
	UserID   string
}

// expirationTime is the time in seconds until the token expires.
func GenerateJWTToken(api *types.Api, data JWTField, secretKey []byte, expirationTime int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": data.UserName,
		"iss": api.AppName,
		"aud": api.AppName,
		"id":  data.UserID,
		"exp": expirationTime,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateJWTToken(tokenString string, secretKey []byte) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	}, jwt.WithValidMethods([]string{"HS256"}))

	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return username, nil
	}

	return "", fmt.Errorf("invalid token")
}
