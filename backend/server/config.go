package server

import "os"

var AllowedPaths = []string{
	// Backend
	"api",
	// Frontend
	"favicon.ico",
	"_nuxt",
	".well-known",
	"200.html",
	"404.html",
	"robots.txt",
}

type Config struct {
	Host     string
	Port     string
	DBUrl    string
	RedisUrl string
}

func LoadEnvConfig() *Config {
	port, is_port_set := os.LookupEnv("GO_PORT")

	if !is_port_set {
		panic("GO_PORT variable is to be set")
	}

	host, is_host_set := os.LookupEnv("GO_HOST")

	if !is_host_set {
		panic("GO_HOST variable is to be set")
	}

	db_url, is_url_set := os.LookupEnv("POSTGRES_URL")

	if !is_url_set {
		panic("POSTGRES_URL variable is to be set")
	}

	redis_url, is_url_set := os.LookupEnv("REDIS_URL")

	if !is_url_set {
		panic("REDIS_URL variable is to be set")
	}

	return &Config{
		Port:     port,
		Host:     host,
		DBUrl:    db_url,
		RedisUrl: redis_url,
	}
}
