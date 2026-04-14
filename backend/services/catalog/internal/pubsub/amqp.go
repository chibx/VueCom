package pubsub

import (
	"os"

	"github.com/chibx/vuecom/backend/services/catalog/internal/global"
	"github.com/chibx/vuecom/backend/services/catalog/internal/utils"
	"github.com/chibx/vuecom/backend/shared/events"
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

func InitPubSub() {
	if DefPubSub != nil {
		return
	}
	mqUrl := getEnv("APP_RABBITMQ_URL")
	DefPubSub, err := NewPubSub(mqUrl)
	utils.FailOnError(err, "Could not connect to rabbitmq instance")

	for _, v := range events.Queues {
		err = DefPubSub.CreateQueue(v)
		utils.FailOnError(err, "Failed to create queue")
	}
	global.Logger.Info("RabbitMQ instance connected")
}
