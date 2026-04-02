package pubsub

import (
	"context"
	"os"
	"time"

	"github.com/chibx/vuecom/backend/services/orders/internal/global"
	"github.com/chibx/vuecom/backend/services/orders/internal/utils"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
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

func InitPubSub() (*amqp.Connection, *amqp.Channel) {
	mqUrl := getEnv("APP_RABBITMQ_URL")
	queueName := getEnv("APP_RABBITMQ_QUEUE")
	conn, err := amqp.Dial(mqUrl)
	utils.FailOnError(err, "Could not connect to rabbitmq instance")
	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		queueName, // name
		true,      // durability
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		amqp.Table{
			amqp.QueueTypeArg: amqp.QueueTypeQuorum,
		},
	)
	utils.FailOnError(err, "Failed to declare a queue")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	body := "Hello World!"
	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	utils.FailOnError(err, "Failed to publish a message")
	global.Logger.Info(" [x] Sent %s\n", zap.String("body", body))

	return conn, ch
}
