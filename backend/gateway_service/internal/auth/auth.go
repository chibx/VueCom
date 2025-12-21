package auth

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRefreshToken (as before).
func GenerateRefreshToken() (string, error) {
	// ... (crypto/rand + hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
