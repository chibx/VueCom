package main

import (
	"context"
	v1 "vuecom/gateway/api/v1"
	"vuecom/gateway/config"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func loadCloudinary() *cloudinary.Cloudinary {
	cldKey := config.GetEnv("CLOUDINARY_KEY")
	cldSecret := config.GetEnv("CLOUDINARY_SECRET")
	cldName := config.GetEnv("CLOUDINARY_CLOUD_NAME")
	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)

	if err != nil {
		panic("Error setting up Cloudinary!!!")
	}

	return cld
}

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
