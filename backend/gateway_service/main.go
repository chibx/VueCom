package main

import (
	"fmt"

	// "sync"
	v1 "vuecom/gateway/api/v1"
	"vuecom/gateway/config"
	"vuecom/gateway/internal/v1/types"
	"vuecom/shared/deps"

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
	v1_api := &types.Api{Config: config, Deps: &deps.Deps{}}

	plugDB(v1_api)
	plugRedis(v1_api)
	plugCloudinary(v1_api)

	err := migrate(v1_api.Deps.DB)
	if err != nil {
		panic("Error while migration")
	}

	v1.LoadRoutes(app, v1_api)

	app.Static("/", "./dist")

	app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
}
