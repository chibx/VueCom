package pubsub

import (
	"context"
	"fmt"

	"github.com/chibx/vuecom/backend/services/payment/internal/global"
	"github.com/goccy/go-json"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

type Handler func(json.RawMessage) error

type Message struct {
	Event   string          `json:"event"`
	Payload json.RawMessage `json:"payload"`
}

type PubSub struct {
	started  bool
	conn     *amqp.Connection
	channel  *amqp.Channel
	handlers map[string]Handler
}

var DefPubSub *PubSub

func NewPubSub(amqpURL string) (*PubSub, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &PubSub{
		conn:     conn,
		channel:  ch,
		handlers: make(map[string]Handler),
	}, nil
}

func (p *PubSub) CreateQueue(queueName string) error {
	_, err := p.channel.QueueDeclare(queueName, true, false, false, false, nil)
	return err
}

// RegisterHandler lets you add new event handlers from anywhere
func (c *PubSub) RegisterHandler(event string, handler Handler) {
	c.handlers[event] = handler
}

func (c *PubSub) Start(ctx context.Context, queueName string) error {
	if c.started {
		global.Logger.Warn("Already started queue consumer")
		return nil
	}

	msgs, err := c.channel.Consume(
		queueName,
		"",    // consumer tag
		false, // autoAck = false
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			handler, ok := c.handlers[msg.Type]
			if !ok {
				global.Logger.Warn("no handler for event", zap.String("event", msg.Type))
				msg.Nack(false, false) // do not requeue
				continue
			}

			if err := handler(msg.Body); err != nil {
				global.Logger.Error(fmt.Sprintf("handler for %s failed", msg.Type), zap.Error(err))
				// You can decide here: requeue, dead-letter, or drop
				msg.Nack(false, true)
				continue
			}

			msg.Ack(false)
		}
	}()
	c.started = true
	// <-ctx.Done() // graceful shutdown
	return nil
}

func (p *PubSub) Publish(queueName string, event string, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		global.Logger.Error("Failed to serialize message payload", zap.Error(err))
		return err
	}

	return p.channel.Publish("", queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Type:        event,
	})
}

func (c *PubSub) Close() {
	c.channel.Close()
	c.conn.Close()
}
