package types

import "vuecom/shared/deps"

type Config struct {
	Host string
	Port string
	// PostgresDSN   string
	// RedisUrl      string
	AllowedPaths  []string
	MockAdminSlug string
	ApiEncKey     []byte // For the API Keys
	SecretKey     []byte // For encrypting jwt
	DbEncKey      []byte // For encrypting db credentials like user information i.e address and password (after hashing of course)
	IsSaas        bool
}

type Api struct {
	// DB       *gorm.DB
	// Redis    *redis.Client
	// Cld      *cloudinary.Cloudinary
	Deps      *deps.Deps
	Config    *Config
	HasAdmin  bool
	IsAppInit bool
	AppName   string
}
