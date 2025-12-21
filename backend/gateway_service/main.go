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
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	config := config.GetConfig()
	v1_api := &types.Api{Config: config, Deps: &deps.Deps{}}

	initServer(v1_api)

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		JSONEncoder:           json.Marshal,
		JSONDecoder:           json.Unmarshal,
		StreamRequestBody:     true,
	})

	app.Use(helmet.New())
	// TODO: Add a rate limiter middleware
	// app.Get("/", func(ctx *fiber.Ctx) error {
	// 	return ctx.SendStatus(fiber.StatusNotFound)
	// })

	v1.LoadRoutes(app, v1_api)

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port)))
}
