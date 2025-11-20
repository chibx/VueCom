package config

type Config struct {
	Host          string
	Port          string
	PostgresDSN   string
	RedisUrl      string
	AllowedPaths  []string
	MockAdminSlug string
	ApiMasterKey  []byte
	isSaas        bool
}

func GetConfig() *Config {
	return &Config{
		Host:          getEnv("GO_HOST", "127.0.0.1"),
		Port:          getEnv("GO_PORT", "2500"),
		PostgresDSN:   loadPostgresDSN(),
		RedisUrl:      getEnv("REDIS_URL"),
		AllowedPaths:  allowedPaths,
		MockAdminSlug: "admin123",
		ApiMasterKey:  loadMasterKey(),
		isSaas:        isSaaS(),
	}
}
