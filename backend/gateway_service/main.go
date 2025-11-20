package main

import (
	"context"
	"fmt"

	// "sync"
	v1 "vuecom/gateway/api/v1"
	"vuecom/gateway/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	_ "github.com/joho/godotenv/autoload"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func plugDB(api *v1.Api, dsn string) {
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}

	api.DB = db
}

func plugRedis(api *v1.Api, redisUrl string) {
	// db, err := gorm.Open(postgres.Open(dsn))

	// if err != nil {
	// 	panic(err)
	// }

	// api.DB = db
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("REDIS_URL should be set!!!")
	}

	client := redis.NewClient(opts)
	cmd := client.Ping(context.Background())
	if cmd.Err() != nil {
		panic("Could not connect to Redis!!!")
	}
	api.Redis = client
}

func main() {
	config := config.GetConfig()

	app := fiber.New(fiber.Config{DisableStartupMessage: true})

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
