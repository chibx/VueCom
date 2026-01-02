package main

import (
	"fmt"
	"os"
	"time"
	"vuecom/gateway/internal/db/gorm_pg"

	_ "github.com/joho/godotenv/autoload"

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
	host := getEnv("PG_HOST")

	user := getEnv("PG_USER")

	passwd := getEnv("PG_PASSWD")

	dbName := getEnv("PG_DBNAME")

	port := getEnv("PG_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", host, user, passwd, dbName, port)
}

func main() {
	now := time.Now()
	dsn := loadPostgresDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	err = gorm_pg.NewGormPGDatabase(db).Migrate()

	if err != nil {
		fmt.Println("failed to migrate: " + err.Error())
		return
	}

	fmt.Println("Migration successful")
	fmt.Println("Migration done in", time.Since(now).Milliseconds(), "ms")
}
