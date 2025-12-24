package types

import (
	"vuecom/shared/models"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/redis/go-redis/v9"
)

// "gorm.io/gorm"

type Deps struct {
	// DB    *gorm.DB
	DB    Database
	Redis *redis.Client
	Cld   *cloudinary.Cloudinary
}

type Config struct {
	Host string
	Port string
	// PostgresDSN   string
	// RedisUrl      string
	ApiEncKey    []byte // For the API Keys
	SecretKey    []byte // For encrypting jwt
	DbEncKey     []byte // For encrypting db credentials like user information i.e address and password (after hashing of course)
	IsSaas       bool
	AllowedPaths []string
}

type Api struct {
	// DB       *gorm.DB
	// Redis    *redis.Client
	// Cld      *cloudinary.Cloudinary
	Deps        *Deps
	Config      *Config
	HasAdmin    bool
	IsAppInit   bool
	AppName     string
	AdminSlug   string
	AppSettings models.AppSettings
}
