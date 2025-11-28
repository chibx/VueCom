package auth

import (
	"crypto/rand"
	"encoding/base64"

	"crypto/aes"
	"crypto/cipher"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// ...

/*
*
* Can be generated with:
package main

import (

	"crypto/rand"
	"encoding/base64"
	"fmt"

)

	func main() {
	    key := make([]byte, 32)
	    if _, err := rand.Read(key); err != nil {
	        panic(err)
	    }
	    fmt.Println("Master Key (base64):", base64.StdEncoding.EncodeToString(key))
	}

*
*/

func GenerateHashFromPass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 0)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func VerifyHashWithPass(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	return err == nil
}

func GenerateAPIKey(length int) (string, error) {
	bytes := make([]byte, length) // e.g., 32 for ~256 bits entropy
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// Encrypt encrypts plaintext with the master key (AES-256-GCM).
// Returns base64-encoded ciphertext (nonce + tag + encrypted data) for easy storage.
func EncryptAPIKey(plaintext string, masterKey []byte) (string, error) {
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nonce, nonce, []byte(plaintext), nil) // Nonce prepended
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts base64-encoded ciphertext back to plaintext.
func DecryptAPIKey(ciphertextBase64 string, masterKey []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesgcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
