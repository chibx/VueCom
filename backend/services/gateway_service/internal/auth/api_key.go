package auth

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateAPIKey(length int) (string, error) {
	bytes := make([]byte, length) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
