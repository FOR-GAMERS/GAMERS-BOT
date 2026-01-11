package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/handlers"
	amqp "github.com/rabbitmq/amqp091-go"
)

// Consumer consumes request messages from RabbitMQ
type Consumer struct {
	channel       *amqp.Channel
	queueName     string
	prefetchCount int
	handlers      map[EventType]handlers.Handler
	bot           *bot.DiscordBot
	publisher     *Publisher
}

// NewConsumer creates a new Consumer
func NewConsumer(conn *amqp.Connection, queueName string, prefetchCount int, bot *bot.DiscordBot, publisher *Publisher) (*Consumer, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Set QoS (prefetch count)
	err = channel.Qos(
		prefetchCount, // prefetch count
		0,             // prefetch size
		false,         // global
	)
	if err != nil {
		channel.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	// Declare the request queue (durable)
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Consumer{
		channel:       channel,
		queueName:     queueName,
		prefetchCount: prefetchCount,
		handlers:      make(map[EventType]handlers.Handler),
		bot:           bot,
		publisher:     publisher,
	}, nil
}

// RegisterHandler registers a handler for a specific event type
func (c *Consumer) RegisterHandler(eventType EventType, handler handlers.Handler) {
	c.handlers[eventType] = handler
}

// Start starts consuming messages
func (c *Consumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queueName, // queue
		"",          // consumer
		false,       // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	slog.Info("Started consuming messages", "queue", c.queueName)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Consumer stopped")
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Message channel closed")
				return fmt.Errorf("message channel closed")
			}
			c.handleMessage(ctx, msg)
		}
	}
}

// handleMessage processes a single message
func (c *Consumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	// Parse request message
	var request RequestMessage
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		slog.Error("Failed to unmarshal request", "error", err)
		c.sendErrorResponse(ctx, "", fmt.Errorf("invalid message format: %w", err))
		msg.Nack(false, false) // Don't requeue invalid messages
		return
	}

	slog.Info("Processing event", "correlation_id", request.CorrelationID, "guild_id", request.GuildID, "event_type", request.EventType)

	// Validate guild_id
	if request.GuildID == "" {
		slog.Error("Missing guild_id in request")
		c.sendErrorResponse(ctx, request.CorrelationID, fmt.Errorf("guild_id is required"))
		msg.Nack(false, false) // Don't requeue invalid messages
		return
	}

	// Get handler for event type
	handler, ok := c.handlers[request.EventType]
	if !ok {
		slog.Error("Unsupported event type", "event_type", request.EventType)
		c.sendErrorResponse(ctx, request.CorrelationID, fmt.Errorf("unsupported event type: %s", request.EventType))
		msg.Nack(false, false) // Don't requeue unsupported events
		return
	}

	// Handle the event
	data, err := handler.Handle(ctx, c.bot, request.GuildID, request.Payload)
	if err != nil {
		slog.Error("Handler failed", "correlation_id", request.CorrelationID, "error", err)

		// Check if error is retriable
		if isRetriableError(err) {
			slog.Info("Requeuing message", "correlation_id", request.CorrelationID)
			msg.Nack(false, true) // Requeue for retry
			return
		}

		// Non-retriable error: send error response and don't requeue
		c.sendErrorResponse(ctx, request.CorrelationID, err)
		msg.Nack(false, false)
		return
	}

	// Send success response
	c.sendSuccessResponse(ctx, request.CorrelationID, data)
	msg.Ack(false)
	slog.Info("Event processed successfully", "correlation_id", request.CorrelationID)
}

// sendSuccessResponse sends a success response
func (c *Consumer) sendSuccessResponse(ctx context.Context, correlationID string, data map[string]interface{}) {
	response := &ResponseMessage{
		CorrelationID: correlationID,
		Success:       true,
		Data:          data,
	}

	if err := c.publisher.Publish(ctx, response); err != nil {
		slog.Error("Failed to publish success response", "correlation_id", correlationID, "error", err)
	}
}

// sendErrorResponse sends an error response
func (c *Consumer) sendErrorResponse(ctx context.Context, correlationID string, err error) {
	response := &ResponseMessage{
		CorrelationID: correlationID,
		Success:       false,
		Error:         err.Error(),
	}

	if err := c.publisher.Publish(ctx, response); err != nil {
		slog.Error("Failed to publish error response", "correlation_id", correlationID, "error", err)
	}
}

// isRetriableError determines if an error is retriable
func isRetriableError(err error) bool {
	// For now, we consider most errors as non-retriable
	// In the future, you can add logic to check for specific error types
	// that indicate temporary issues (network timeouts, rate limits, etc.)
	return false
}

// Close closes the consumer channel
func (c *Consumer) Close() error {
	if c.channel != nil {
		return c.channel.Close()
	}
	return nil
}
