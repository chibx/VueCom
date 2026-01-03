package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
	"vuecom/gateway/config"
	"vuecom/gateway/internal/db/gorm_pg"
	"vuecom/gateway/internal/types"
	dbModels "vuecom/shared/models/db"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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

	api.Deps.DB = gorm_pg.NewGormPGDatabase(db)
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

// func attachSentry(app *fiber.App) {
// 	sentry_dsn := getEnv("SENTRY_DSN", "")
// 	if sentry_dsn != "" {
// 		if err := sentry.Init(sentry.ClientOptions{
// 			Dsn: getEnv("SENTRY_DSN"),
// 		}); err != nil {
// 			fmt.Printf("Sentry initialization failed: %v\n", err)
// 		}

// 		// Later in the code
// 		sentryHandler := sentryfiber.New(sentryfiber.Options{
// 			WaitForDelivery: true,
// 		})

// 		app.Use(sentryHandler)
// 	} else {
// 		fmt.Println("Skipping Sentry Initialization! SENTRY_DSN not found")
// 	}
// }

func appIfInitialized(api *types.Api) (*dbModels.AppData, error) {
	logger := api.Deps.Logger
	appData, err := api.Deps.DB.AppData().GetAppData(context.Background())

	if err != nil {
		if errors.Is(err, types.ErrDbNil) {
			logger.Info("No active app found in DB")
			return &dbModels.AppData{}, err
		}
		logger.Error("Error occurerd while fetching app data", zap.Error(err))
		return nil, err
	}

	logger.Info("App data fetched successfully from DB", zap.String("name", appData.Name))
	return appData, nil
}

func checkIfOwnerExists(api *types.Api) (bool, error) {
	logger := api.Deps.Logger
	count, err := api.Deps.DB.AppData().CountOwner(context.Background())

	if err != nil {
		logger.Error("Error occurerd while checking for owner existence", zap.Error(err))
		return false, err
	}

	logger.Info("Owner existence check complete", zap.Int64("owner_count", count))
	return count > 0, nil
}

func initLogger(v1_api *types.Api) {
	writer := zapcore.AddSync(os.Stdout) // Use standard output as the log target
	zapPreset := zap.NewProductionEncoderConfig()
	zapPreset.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(zapPreset),
		writer,
		zapcore.InfoLevel,
	)

	core = zapcore.NewSamplerWithOptions(core, time.Second, 10, 5)

	logger := zap.New(core)

	v1_api.Deps.Logger = logger
}

func initServer(_ *fiber.App, v1_api *types.Api) {
	initLogger(v1_api)
	plugDB(v1_api)
	plugRedis(v1_api)
	plugCloudinary(v1_api)
	// attachSentry(app)
	logger := v1_api.Deps.Logger
	now := time.Now()
	// err := migrate(v1_api.Deps.DB)
	err := v1_api.Deps.DB.Migrate()
	if err != nil {
		panic("Error while migration")
	}
	logger.Info("Auto Migration took", zap.String("duration", strconv.Itoa(int(time.Since(now).Milliseconds()))+"ms"))

	appData, _ := appIfInitialized(v1_api)
	v1_api.HasAdmin, _ = checkIfOwnerExists(v1_api)

	v1_api.IsAppInit = appData.Name != ""

	if len(v1_api.AppName) > 0 {
		v1_api.AppName = appData.Name
	} else {
		logger.Warn("App Name not found in DB, using default 'Vuecom_test'")
		v1_api.AppName = "Vuecom_test"
	}

	if len(appData.AdminRoute) > 0 {
		v1_api.AdminSlug = appData.AdminRoute
	} else {
		v1_api.AdminSlug = "admin123"
		logger.Warn("Admin Route not found in DB, using default 'admin123'")
	}
}
