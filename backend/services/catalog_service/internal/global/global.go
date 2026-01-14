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
	host := getEnv("CATALOG_PG_HOST")
	user := getEnv("CATALOG_PG_USER")
	passwd := getEnv("CATALOG_PG_PASSWD")
	dbName := getEnv("CATALOG_PG_DBNAME")
	port := getEnv("CATALOG_PG_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, passwd, dbName, port)
}

func newDB() *gorm.DB {
	dsn := loadPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn))

	if err != nil {
		panic(err)
	}
	return db
}

func newCloudinary() *cloudinary.Cloudinary {
	cldKey := getEnv("CLOUDINARY_KEY")
	cldSecret := getEnv("CLOUDINARY_SECRET")
	cldName := getEnv("CLOUDINARY_CLOUD_NAME")
	cld, err := cloudinary.NewFromParams(cldName, cldKey, cldSecret)

	if err != nil {
		panic("Error setting up Cloudinary!!!")
	}

	return cld
}

func newRedis() *redis.Client {
	redisUrl := getEnv("CATALOG_REDIS_URL")
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("CATALOG_REDIS_URL should be set!!!")
	}

	client := redis.NewClient(opts)
	cmd := client.Ping(context.Background())
	if cmd.Err() != nil {
		panic("Could not connect to Redis!!!")
	}

	return client
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

var (
	Repo   = db.NewCatalogDB(newDB())
	Logger = newLogger()
	Redis  = newRedis()
	Cld    = newCloudinary()
)
