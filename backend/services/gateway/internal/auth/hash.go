package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"errors"

	"golang.org/x/crypto/argon2"
)

type HashParams struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultHashParams = &HashParams{
	Memory:      64 * 1024,
	Iterations:  2,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   32,
}

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

// error messages
var (
	ErrInvalidHash         = errors.New("the encoded hash is not in the correct format")
	ErrIncompatibleVersion = errors.New("incompatible version of argon2")
)

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateHashFromString(password string, p *HashParams) (encodedHash string, err error) {
	salt, err := generateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), []byte(salt), p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	//we base64 encoded both the salt and the hash.
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)
	base64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.Memory, p.Iterations, p.Parallelism, base64Salt, base64Hash)

	return encodedHash, nil
}

func splitHash(encodedHash string) (p *HashParams, salt, hash []byte, err error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return nil, nil, nil, ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(parts[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}

	if version != argon2.Version {
		return nil, nil, nil, ErrIncompatibleVersion
	}

	p = &HashParams{}
	_, err = fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.Memory, &p.Iterations, &p.Parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(parts[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.SaltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(parts[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.KeyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func CompareRawAndHash(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := splitHash(encodedHash)
	if err != nil {
		return false, err
	}
	//reproduced from the password from login form
	inputHash := argon2.IDKey([]byte(password), []byte(salt), p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	//we use ConstantTimeCompare to help mitigate timing attacks.
	if subtle.ConstantTimeCompare(hash, inputHash) == 1 {
		return true, nil
	}

	return false, nil
}

func Encrypt(text string, secretKey []byte) (encryptedHash string, err error) {
	//setting up a cipher block from the aes.NewCipher method. We use a 32 bytes key here.
	cipherBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}
	//Create a random nonce with the appropriate size.
	nonce, err := generateRandomBytes(uint32(gcm.NonceSize()))
	if err != nil {
		return "", err
	}
	//Use the Seal method with the nonce to encrypt our data.
	cipherText := gcm.Seal(nonce, nonce, []byte(text), nil)

	//The Seal method returns a byte slice, hence we encode it with base64 (or hex) to store in database.
	encryptedHash = base64.RawStdEncoding.EncodeToString(cipherText)

	return encryptedHash, nil
}

func Decrypt(encryptedHash string, secretKey []byte) (resultText string, err error) {
	cipherText, err := base64.RawStdEncoding.DecodeString(encryptedHash)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]
	dec, err := gcm.Open(nil, nonce, cipherText, nil)
	resultText = string(dec)
	if err != nil {
		return "", err
	}

	return string(resultText), nil
}
