package types

import (
	"github.com/chibx/vuecom/backend/shared/models"

	"github.com/chibx/vuecom/backend/services/gateway/internal/db/gorm_pg"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
)

// "gorm.io/gorm"

type Deps struct {
	// DB    *gorm.DB
	DB      *gorm_pg.GormPGDatabase
	Redis   *redis.Client
	Cld     *cloudinary.Cloudinary
	Limiter *redis_rate.Limiter
}

type Config struct {
	Host string
	Port string
	// PostgresDSN   string
	// RedisUrl      string
	ApiEncKey []byte // For the API Keys
	// SecretKey    []byte // For encrypting jwt
	SecretKey    []byte // For encrypting db credentials like user information i.e address and password (after hashing of course)
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
