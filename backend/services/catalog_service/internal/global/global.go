package global

import (
	"context"
	"fmt"
	"os"
	"time"
	"vuecom/catalog/internal/db"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Logger = newLogger()

var (
	Repo  = db.NewCatalogDB(newDB())
	Redis = newRedis()
	Cld   = newCloudinary()
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
	dbName := getEnv("CATALOG_PG_DBNAME")
	port := getEnv("APP_PG_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, passwd, dbName, port)
}

func newDB() *gorm.DB {
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
			Logger.Error("failed to initialize db conn for catalog service", zap.Error(err))
		} else {
			return db
		}

		Logger.Info("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	Logger.Panic("Could not connect to database after multiple retries")
	return nil
}

func newCloudinary() *cloudinary.Cloudinary {
	cldKey := getEnv("CLOUDINARY_KEY")
	cldSecret := getEnv("CLOUDINARY_SECRET")
	cldName := getEnv("CLOUDINARY_CLOUD_NAME")
	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)

	if err != nil {
		Logger.Error("failed to initialize cloudinary for catalog service", zap.Error(err))
		Logger.Panic("Error setting up Cloudinary!!!")
	}

	return cld
}

func newRedis() *redis.Client {
	redisUrl := getEnv("APP_REDIS_URL")
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		Logger.Error("failed to parse redis url for catalog service", zap.Error(err))
		Logger.Panic("APP_REDIS_URL should be set!!!")
	}

	client := redis.NewClient(opts)

	for range 5 {
		cmd := client.Ping(context.Background())
		err = cmd.Err()
		if err != nil {
			Logger.Error("failed to connect to redis for catalog service", zap.Error(err))
		} else {
			return client
		}

		Logger.Info("Backing off for 2 seconds...")
		time.Sleep(2 * time.Second)
	}

	Logger.Panic("Could not connect to Redis after multiple retries")
	return nil
}

func newLogger() *zap.Logger {
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

	return logger
}
