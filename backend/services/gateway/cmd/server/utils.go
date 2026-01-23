package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	// "strconv"
	"time"

	appModels "github.com/chibx/vuecom/backend/shared/models/db/appdata"

	"github.com/chibx/vuecom/backend/services/gateway/config"
	"github.com/chibx/vuecom/backend/services/gateway/internal/db/gorm_pg"
	"github.com/chibx/vuecom/backend/services/gateway/internal/types"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/go-redis/redis_rate/v10"
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

	host := getEnv("APP_PG_HOST")

	user := getEnv("APP_PG_USER")

	passwd := getEnv("APP_PG_PASSWORD")

	dbName := getEnv("GATE_PG_DBNAME")

	port := getEnv("APP_PG_PORT")

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
	logger := api.Deps.Logger
	dsn := loadPostgresDSN()
	var db *gorm.DB
	var err error
	// if err != nil {
	// 	logger.Error("failed to initialize db conn", zap.Error(err))
	// 	panic(err)
	// }

	for range 5 {
		db, err = gorm.Open(postgres.Open(dsn))

		if err != nil {
			logger.Info("Could not open postgres connection", zap.Error(err))
		} else {
			api.Deps.DB = gorm_pg.NewGormPGDatabase(db)
			return
		}

		logger.Info("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	logger.Panic("Could not connect to database after multiple retries")
}

func plugRedis(api *types.Api) {
	logger := api.Deps.Logger
	redisUrl := getEnv("APP_REDIS_URL")
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		logger.Error("failed to parse redis url", zap.Error(err))
		panic("APP_REDIS_URL should be set!!!")
	}

	client := redis.NewClient(opts)

	for range 5 {
		cmd := client.Ping(context.Background())
		err = cmd.Err()
		if err != nil {
			logger.Error("failed to connect to redis", zap.Error(err))
		} else {
			api.Deps.Redis = client
			return
		}

		logger.Info("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}
}

func setupLimiter(api *types.Api) {
	rdb := api.Deps.Redis
	api.Deps.Limiter = redis_rate.NewLimiter(rdb)
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

func appIfInitialized(api *types.Api) (*appModels.AppData, error) {
	logger := api.Deps.Logger
	appData, err := api.Deps.DB.AppData().GetAppData(context.Background())

	if err != nil {
		if errors.Is(err, types.ErrDbNil) {
			logger.Info("No active app found in DB")
			return &appModels.AppData{}, err
		}
		logger.Error("Error occurerd while fetching app data", zap.Error(err))
		return nil, err
	}

	logger.Info("App data fetched successfully from DB", zap.String("name", appData.AppName))
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
	setupLimiter(v1_api)
	plugCloudinary(v1_api)
	// attachSentry(app)

	// ---------------------------
	// appEnv := os.Getenv("APP_ENVIRONMENT")
	// if appEnv != "production" {
	// 	logger := v1_api.Deps.Logger
	// 	// Migrate DB
	// 	now := time.Now()
	// 	err := v1_api.Deps.DB.Migrate()
	// 	if err != nil {
	// 		panic("Error while migration")
	// 	}
	// 	logger.Info("Auto Migration took", zap.String("duration", strconv.Itoa(int(time.Since(now).Milliseconds()))+"ms"))
	// }
	// --------------------------

	appData, _ := appIfInitialized(v1_api)
	v1_api.HasAdmin, _ = checkIfOwnerExists(v1_api)

	if appData != nil {
		v1_api.IsAppInit = appData.AppName != ""

		if len(appData.AppName) > 0 {
			v1_api.AppName = appData.AppName
		} /* else {
			logger.Warn("App Name not found in DB, using default 'Vuecom_test'")
			v1_api.AppName = "Vuecom_test"
		} */

		if len(appData.AdminRoute) > 0 {
			v1_api.AdminSlug = appData.AdminRoute
		} /*  else {
			v1_api.AdminSlug = "admin123"
			logger.Warn("Admin Route not found in DB, using default 'admin123'")
		} */
	}
}
