package config

import "vuecom/gateway/internal/v1/types"

// type Config struct {
// 	Host string
// 	Port string
// 	// PostgresDSN   string
// 	// RedisUrl      string
// 	AllowedPaths  []string
// 	MockAdminSlug string
// 	ApiMasterKey  []byte
// 	isSaas        bool
// }

func GetConfig() *types.Config {
	return &types.Config{
		Host: GetEnv("GO_HOST", "127.0.0.1"),
		Port: GetEnv("GO_PORT", "2500"),
		// PostgresDSN:   loadPostgresDSN(),
		// RedisUrl:      GetEnv("REDIS_URL"),
		AllowedPaths: allowedPaths,
		ApiEncKey:    loadKey("API_ENC_KEY"),
		SecretKey:    loadKey("SECRET_KEY"),
		DbEncKey:     loadKey("DB_ENC_KEY"),
		IsSaas:       isSaaS(),
	}
}
