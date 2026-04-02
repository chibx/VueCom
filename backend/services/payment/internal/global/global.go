package global

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/chibx/vuecom/backend/services/payment/internal/db"
	"github.com/chibx/vuecom/backend/shared/types"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Logger = newLogger("[Payment]: ")

var (
	Repo  = db.NewPaymentDB(newDB())
	Redis = newRedis()
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
	dbName := getEnv("PAYMENT_PG_DBNAME")
	port := getEnv("APP_PG_PORT")

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

func newRedis() *redis.Client {
	redisUrl := getEnv("APP_REDIS_URL")
	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		panic("APP_REDIS_URL should be set!!!")
	}

	client := redis.NewClient(opts)
	cmd := client.Ping(context.Background())
	if cmd.Err() != nil {
		panic("Could not connect to Redis!!!")
	}

	return client
}

func newLogger(prefix string) *zap.Logger {
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

	logger = logger.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return types.NewZapPrefix(c, prefix)
	}))

	return logger
}
