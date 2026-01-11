package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher publishes response messages to RabbitMQ
type Publisher struct {
	channel   *amqp.Channel
	queueName string
}

// NewPublisher creates a new Publisher
func NewPublisher(conn *amqp.Connection, queueName string) (*Publisher, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare the response queue (durable)
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		err := channel.Close()
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Publisher{
		channel:   channel,
		queueName: queueName,
	}, nil
}

// Publish publishes a response message to the response queue
func (p *Publisher) Publish(ctx context.Context, response *ResponseMessage) error {
	body, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		"",          // exchange
		p.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	slog.Debug("Published response", "correlation_id", response.CorrelationID, "success", response.Success)
	return nil
}

// Close closes the publisher channel
func (p *Publisher) Close() error {
	if p.channel != nil {
		return p.channel.Close()
	}
	return nil
}
