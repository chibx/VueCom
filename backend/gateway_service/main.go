package main

import (
	"fmt"

	// "sync"
	v1 "vuecom/gateway/api/v1"
	"vuecom/gateway/config"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := config.GetConfig()

	app := fiber.New(fiber.Config{DisableStartupMessage: true, JSONEncoder: json.Marshal, JSONDecoder: json.Unmarshal})

	app.Use(helmet.New())
	//! TODO Add a rate limiter middleware
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})
	v1_api := &v1.Api{Config: config}

	plugDB(v1_api, config.PostgresDSN)

	err := migrate(v1_api.DB)
	if err != nil {
		panic("Error while migration")
	}

	v1_api.LoadApis(app)

	app.Static("/", "./dist")

	app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
}
