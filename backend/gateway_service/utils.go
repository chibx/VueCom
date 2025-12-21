package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"
	"vuecom/gateway/config"
	"vuecom/gateway/internal/v1/types"
	dbModels "vuecom/shared/models/db"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getEnv(env string, sub ...string) string {
	val := os.Getenv(env)
	if val == "" {
		if len(sub) > 0 {
			return sub[0]
		}
		panic("Environment Variable " + env + " not set")
	}
	return val
}

func loadPostgresDSN() string {

	// "host=localhost user=gorm password=gorm dbname=gorm port=5432 sslmode=disable"

	host := getEnv("PG_HOST")

	user := getEnv("PG_USER")

	passwd := getEnv("PG_PASSWD")

	dbName := getEnv("PG_DBNAME")

	port := getEnv("PG_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, passwd, dbName, port)
}

func plugCloudinary(api *types.Api) {
	cldKey := config.GetEnv("CLOUDINARY_KEY")
	cldSecret := config.GetEnv("CLOUDINARY_SECRET")
	cldName := config.GetEnv("CLOUDINARY_CLOUD_NAME")
	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)

	if err != nil {
		panic("Error setting up Cloudinary!!!")
	}

	api.Deps.Cld = cld
}

func plugDB(api *types.Api) {
	dsn := loadPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}

	api.Deps.DB = db
}

func plugRedis(api *types.Api) {
	redisUrl := getEnv("REDIS_URL")
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("REDIS_URL should be set!!!")
	}

	client := redis.NewClient(opts)
	cmd := client.Ping(context.Background())
	if cmd.Err() != nil {
		panic("Could not connect to Redis!!!")
	}
	api.Deps.Redis = client
}

func appIfInitialized(api *types.Api) (*dbModels.AppData, error) {
	var appData = &dbModels.AppData{}

	err := api.Deps.DB.First(appData).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dbModels.AppData{}, err
		}
		return nil, err
	}

	return appData, nil
}

func checkIfOwnerExists(api *types.Api) (bool, error) {
	var count int64

	err := api.Deps.DB.Model(&dbModels.BackendUser{}).Where("role = 'owner'").Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

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

	v1_api.IsAppInit = appData.Name != ""

	if len(v1_api.AppName) > 0 {
		v1_api.AppName = appData.Name
	} else {
		v1_api.AppName = "Vuecom_test"
	}

	if len(appData.AdminRoute) > 0 {
		v1_api.AdminSlug = appData.AdminRoute
	} else {
		v1_api.AdminSlug = "admin123"
	}
}
