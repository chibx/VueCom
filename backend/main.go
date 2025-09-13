package main

import (
	"fmt"

	// "sync"
	"vuecom/server"

	"github.com/gofiber/fiber/v2"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	server.Server
}

func main() {
	var config = server.LoadEnvConfig()

	var server Server = Server{Server: server.Server{}}

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})
	// For validating the admin slug
	app.Use("/:admin/*", server.ValidateSlug)

	api := app.Group("/api")
	server.LoadApis(api)

	app.Static("/", "./dist")

	app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
}
