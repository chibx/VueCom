package config

import (
	"encoding/base64"
	"os"
)

const VERSION = "1.0.0"

func GetEnv(env string, sub ...string) string {
	val := os.Getenv(env)
	if val == "" {
		if len(sub) > 0 {
			return sub[0]
		}
		panic("Environment Variable " + env + " not set")
	}
	return val
}

func loadKey(keyEnv string) []byte {
	keyBase64 := GetEnv(keyEnv)

	var err error
	Key, err := base64.StdEncoding.DecodeString(keyBase64)
	if err != nil || len(Key) != 32 {
		panic("Invalid environment variable " + keyEnv)
	}
	return Key
}

func isSaaS() bool {
	saas := GetEnv("IS_SAAS", "false")

	return saas == "true"
}
