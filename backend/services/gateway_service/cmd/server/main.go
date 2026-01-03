package main

import (
	"fmt"

	// "sync"
	v1 "vuecom/gateway/api/v1"
	"vuecom/gateway/config"
	"vuecom/gateway/internal/types"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	// "go.uber.org/zap"
)

func main() {

	// logger.Info()
	// --------------------------------------------------------
	config := config.GetConfig()
	v1_api := &types.Api{Config: config, Deps: &types.Deps{}}

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

	initServer(app, v1_api)
	logger := v1_api.Deps.Logger
	defer func() {
		_ = logger.Sync()
	}()

	v1.LoadRoutes(app, v1_api)

	logger.Fatal("Error starting server:", zap.Error(app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port))))
}
