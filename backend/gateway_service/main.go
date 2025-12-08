package main

import (
	"fmt"
	"time"

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

func initServer(v1_api *types.Api) {
	plugDB(v1_api)
	plugRedis(v1_api)
	plugCloudinary(v1_api)

	now := time.Now()
	err := migrate(v1_api.Deps.DB)
	fmt.Println("Auto Migration took", time.Since(now).Milliseconds(), "ms")
	if err != nil {
		panic("Error while migration")
	}

	appData, _ := appIfInitialized(v1_api)
	v1_api.HasAdmin, _ = checkIfOwnerExists(v1_api)

	v1_api.IsAppInit = appData != nil
	v1_api.AppName = appData.Name
	if len(appData.AdminRoute) > 0 {
		v1_api.AdminSlug = appData.AdminRoute
	} else {
		v1_api.AdminSlug = "admin123"
	}
}

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
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendStatus(fiber.StatusNotFound)

	})

	v1.LoadRoutes(app, v1_api)

	app.Static("/", "./dist")

	log.Fatal(app.Listen(fmt.Sprintf("%s:%s", config.Host, config.Port)))
}
