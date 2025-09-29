package utils

import (
	"fmt"
	"os"
	"strings"
	"vuecom/config"

	"github.com/gofiber/fiber/v2"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/utils"
)

func ExtractRouteParts(route string) []string {
	route = utils.CopyString(route)
	var routeLength = len(route)
	var routeParts []string

	if routeLength == 1 {
		// It's just "/"
		routeParts = []string{""}
	} else {
		var hasTrailingSlash = string(route[routeLength-1]) == "/"

		if hasTrailingSlash {
			route = route[:len(route)-1]
		}
		// First index is gonna be ""
		routeParts = strings.Split(route, "/")[1:]
	}

	return routeParts
}

func ServeIndex(ctx *fiber.Ctx) error {
	return writeFile(ctx, "./dist/index.html", fiber.MIMETextHTMLCharsetUTF8)
}

func writeFile(ctx *fiber.Ctx, path string, ctype string) error {
	file, err := os.ReadFile(path)

	if err != nil {
		log.Error(err)
		return fiber.ErrNotFound
	}

	header := &ctx.Response().Header
	header.SetContentType(ctype)
	header.SetContentLength(len(file))

	_, err = ctx.Write(file)

	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}

	return nil
}

func LoadEnvConfig() *config.Config {
	port, isSet := os.LookupEnv("GO_PORT")

	if !isSet {
		panic("GO_PORT variable is to be set")
	}

	host, isSet := os.LookupEnv("GO_HOST")

	if !isSet {
		panic("GO_HOST variable is to be set")
	}

	postgresDSN := loadPostgresDSN()

	redisUrl, isSet := os.LookupEnv("REDIS_URL")

	if !isSet {
		panic("REDIS_URL variable is to be set")
	}

	return &config.Config{
		Port:        port,
		Host:        host,
		PostgresDSN: postgresDSN,
		RedisUrl:    redisUrl,
	}
}

func loadPostgresDSN() string {

	// "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"

	host, isSet := os.LookupEnv("PG_HOST")

	if !isSet {
		panic("PG_HOST env is required")
	}

	user, isSet := os.LookupEnv("PG_USER")

	if !isSet {
		panic("PG_USER env is required")
	}

	passwd, isSet := os.LookupEnv("PG_PASSWD")

	if !isSet {
		panic("PG_PASSWD env is required")
	}

	dbName, isSet := os.LookupEnv("PG_DBNAME")

	if !isSet {
		panic("PG_DBNAME env is required")
	}

	port, isSet := os.LookupEnv("PG_PORT")

	if !isSet {
		panic("PG_PORT env is required")
	}

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, passwd, dbName, port)
}
