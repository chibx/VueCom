package main

import (
	"fmt"

	// "sync"
	"vuecom/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	var config = api.LoadEnvConfig()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use(helmet.New())
	//! TODO Add a rate limiter middleware
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})

	LoadApis(app)

	app.Static("/", "./dist")

	app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
}
