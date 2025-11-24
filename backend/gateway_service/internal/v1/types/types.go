package types

import "vuecom/shared/deps"

type Config struct {
	Host string
	Port string
	// PostgresDSN   string
	// RedisUrl      string
	AllowedPaths  []string
	MockAdminSlug string
	ApiMasterKey  []byte
	IsSaas        bool
}

type Api struct {
	// DB       *gorm.DB
	// Redis    *redis.Client
	// Cld      *cloudinary.Cloudinary
	Deps     *deps.Deps
	Config   *Config
	HasAdmin bool
}
