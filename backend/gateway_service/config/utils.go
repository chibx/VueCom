package config

import (
	"encoding/base64"
	"fmt"
	"os"
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

func loadMasterKey() []byte {
	keyBase64 := getEnv("API_MASTER_KEY")

	var err error
	masterKey, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil || len(masterKey) != 32 {
		panic("Invalid API_MASTER_KEY")
	}
	return masterKey
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

func isSaaS() bool {
	saas := getEnv("IS_SAAS", "false")

	if saas == "true" {
		return true
	}

	return false
}
