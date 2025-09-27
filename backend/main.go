package main

import (
	"fmt"

	// "sync"
	"vuecom/api"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func plugDB(api *api.Api, dsn string) {
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}

	api.DB = db
}

func main() {
	config := api.LoadEnvConfig()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Use(helmet.New())
	//! TODO Add a rate limiter middleware
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)
	})
	handler := &api.Api{}

	plugDB(handler, config.PostgresDSN)
	LoadApis(app, handler)

	app.Static("/", "./dist")

	app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))
}
