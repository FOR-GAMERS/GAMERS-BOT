package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/handlers"
	amqp "github.com/rabbitmq/amqp091-go"
)

// ConsumerManager manages multiple queue consumers on a single AMQP connection.
// It declares the gamers.events exchange, binds new queues to it, and also
// maintains the legacy discord.commands queue with its request/response pattern.
type ConsumerManager struct {
	conn          *amqp.Connection
	exchange      string
	prefetchCount int
	bot           *bot.DiscordBot
	publisher     *Publisher // used only by legacy consumer for responses
	handlers      map[EventType]handlers.Handler
	channels      []*amqp.Channel
	dedup         *DedupCache
}

// NewConsumerManager creates a new ConsumerManager.
// exchange is the primary exchange name (e.g. "gamers.events").
func NewConsumerManager(conn *amqp.Connection, exchange string, prefetchCount int, bot *bot.DiscordBot, publisher *Publisher) *ConsumerManager {
	return &ConsumerManager{
		conn:          conn,
		exchange:      exchange,
		prefetchCount: prefetchCount,
		bot:           bot,
		publisher:     publisher,
		handlers:      make(map[EventType]handlers.Handler),
		dedup:         NewDedupCache(1 * time.Hour),
	}
}

// RegisterHandler registers a handler for a specific event type
func (cm *ConsumerManager) RegisterHandler(eventType EventType, handler handlers.Handler) {
	cm.handlers[eventType] = handler
}

// SetupTopology declares the primary exchange and sets up all queue bindings from DefaultQueueBindings.
func (cm *ConsumerManager) SetupTopology() error {
	ch, err := cm.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for topology setup: %w", err)
	}
	defer ch.Close()

	// Declare the primary exchange
	err = ch.ExchangeDeclare(
		cm.exchange, // name
		"topic",     // type
		true,        // durable
		false,       // auto-deleted
		false,       // internal
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange %s: %w", cm.exchange, err)
	}
	slog.Info("Exchange declared", "exchange", cm.exchange)

	// Declare and bind all queues
	for _, qb := range DefaultQueueBindings() {
		_, err = ch.QueueDeclare(
			qb.QueueName, // name
			true,         // durable
			false,        // delete when unused
			false,        // exclusive
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", qb.QueueName, err)
		}

		for _, rk := range qb.RoutingKeys {
			err = ch.QueueBind(
				qb.QueueName, // queue name
				rk,           // routing key
				cm.exchange,  // exchange
				false,        // no-wait
				nil,          // arguments
			)
			if err != nil {
				return fmt.Errorf("failed to bind queue %s with key %s: %w", qb.QueueName, rk, err)
			}
			slog.Info("Queue bound to exchange", "queue", qb.QueueName, "exchange", cm.exchange, "routing_key", rk)
		}
	}

	return nil
}

// SetupLegacyQueue declares legacy exchanges and binds the legacy queue.
func (cm *ConsumerManager) SetupLegacyQueue(queueName string, bindings []LegacyBinding) error {
	ch, err := cm.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for legacy setup: %w", err)
	}
	defer ch.Close()

	// Declare the legacy queue
	_, err = ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare legacy queue %s: %w", queueName, err)
	}

	for _, b := range bindings {
		if b.Exchange == "" {
			continue
		}

		err = ch.ExchangeDeclare(
			b.Exchange, // name
			"topic",    // type
			true,       // durable
			false,      // auto-deleted
			false,      // internal
			false,      // no-wait
			nil,        // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to declare legacy exchange %s: %w", b.Exchange, err)
		}

		err = ch.QueueBind(
			queueName,    // queue name
			b.RoutingKey, // routing key
			b.Exchange,   // exchange
			false,        // no-wait
			nil,          // arguments
		)
		if err != nil {
			return fmt.Errorf("failed to bind legacy queue %s to exchange %s: %w", queueName, b.Exchange, err)
		}

		slog.Info("Legacy queue bound to exchange", "queue", queueName, "exchange", b.Exchange, "routing_key", b.RoutingKey)
	}

	return nil
}

// LegacyBinding represents an exchange and routing key binding for the legacy queue
type LegacyBinding struct {
	Exchange   string
	RoutingKey string
}

// Start starts consuming from all queues. It launches a goroutine per queue and blocks until ctx is cancelled.
func (cm *ConsumerManager) Start(ctx context.Context, legacyQueueName string) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 1)

	// Start notification queue consumers (no response publishing)
	for _, qb := range DefaultQueueBindings() {
		wg.Add(1)
		go func(queueName string) {
			defer wg.Done()
			if err := cm.consumeNotificationQueue(ctx, queueName); err != nil {
				select {
				case errCh <- fmt.Errorf("notification queue %s error: %w", queueName, err):
				default:
				}
			}
		}(qb.QueueName)
	}

	// Start legacy queue consumer (with response publishing)
	if legacyQueueName != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := cm.consumeLegacyQueue(ctx, legacyQueueName); err != nil {
				select {
				case errCh <- fmt.Errorf("legacy queue %s error: %w", legacyQueueName, err):
				default:
				}
			}
		}()
	}

	// Wait for all consumers to finish or first error
	doneCh := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneCh)
	}()

	select {
	case <-doneCh:
		return ctx.Err()
	case err := <-errCh:
		return err
	}
}

// consumeNotificationQueue consumes from a notification queue (Ack/Nack only, no response).
func (cm *ConsumerManager) consumeNotificationQueue(ctx context.Context, queueName string) error {
	ch, err := cm.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for queue %s: %w", queueName, err)
	}
	cm.channels = append(cm.channels, ch)

	err = ch.Qos(cm.prefetchCount, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS for queue %s: %w", queueName, err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer for queue %s: %w", queueName, err)
	}

	slog.Info("Started consuming notification queue", "queue", queueName)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Notification consumer stopped", "queue", queueName)
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Notification message channel closed", "queue", queueName)
				return fmt.Errorf("message channel closed for queue %s", queueName)
			}
			cm.handleNotificationMessage(ctx, msg, queueName)
		}
	}
}

// handleNotificationMessage processes a message from a notification queue.
// Dispatches by AMQP header event_type first, falls back to JSON body event_type.
func (cm *ConsumerManager) handleNotificationMessage(ctx context.Context, msg amqp.Delivery, queueName string) {
	slog.Info("Received notification message", "queue", queueName, "routing_key", msg.RoutingKey, "body", string(msg.Body))

	// Determine event type: prefer AMQP header, fallback to JSON body
	eventType := cm.resolveEventType(msg)
	if eventType == "" {
		slog.Error("Cannot determine event_type", "queue", queueName)
		msg.Nack(false, false) // discard: unrecoverable
		return
	}

	// Parse the full message body as payload
	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Body, &payload); err != nil {
		slog.Error("Failed to unmarshal notification payload", "error", err, "queue", queueName)
		msg.Nack(false, false) // discard: malformed JSON won't fix itself
		return
	}

	// Idempotency check: skip duplicate events based on event_id
	eventID, _ := payload["event_id"].(string)
	if cm.dedup.IsDuplicate(eventID) {
		slog.Warn("Duplicate event detected, skipping", "event_id", eventID, "event_type", eventType)
		msg.Ack(false) // ack so it doesn't get redelivered
		return
	}

	slog.Info("Dispatching notification event", "queue", queueName, "event_type", eventType, "event_id", eventID)

	handler, ok := cm.handlers[eventType]
	if !ok {
		slog.Warn("No handler registered for event type", "event_type", eventType, "queue", queueName)
		msg.Nack(false, false) // discard: no handler will magically appear
		return
	}

	// Extract and validate guild_id
	guildID := extractGuildID(payload)
	if guildID == "" {
		slog.Warn("Missing guild_id in notification event, proceeding without it",
			"event_type", eventType, "event_id", eventID)
	}

	// Handle the event
	_, err := handler.Handle(ctx, cm.bot, guildID, payload)
	if err != nil {
		slog.Error("Notification handler failed",
			"event_type", eventType, "queue", queueName,
			"event_id", eventID, "error", err)
		// Requeue for retry on handler errors (e.g. transient Discord API failures)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
	slog.Info("Notification event processed", "event_type", eventType, "queue", queueName, "event_id", eventID)
}

// consumeLegacyQueue consumes from the legacy queue with request/response pattern.
func (cm *ConsumerManager) consumeLegacyQueue(ctx context.Context, queueName string) error {
	ch, err := cm.conn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open channel for legacy queue %s: %w", queueName, err)
	}
	cm.channels = append(cm.channels, ch)

	err = ch.Qos(cm.prefetchCount, 0, false)
	if err != nil {
		return fmt.Errorf("failed to set QoS for legacy queue %s: %w", queueName, err)
	}

	msgs, err := ch.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer for legacy queue %s: %w", queueName, err)
	}

	slog.Info("Started consuming legacy queue", "queue", queueName)

	for {
		select {
		case <-ctx.Done():
			slog.Info("Legacy consumer stopped", "queue", queueName)
			return ctx.Err()
		case msg, ok := <-msgs:
			if !ok {
				slog.Warn("Legacy message channel closed", "queue", queueName)
				return fmt.Errorf("message channel closed for legacy queue %s", queueName)
			}
			cm.handleLegacyMessage(ctx, msg)
		}
	}
}

// handleLegacyMessage processes a message from the legacy queue (request/response pattern).
func (cm *ConsumerManager) handleLegacyMessage(ctx context.Context, msg amqp.Delivery) {
	slog.Info("Received legacy message", "body", string(msg.Body))

	// Parse request message
	var request RequestMessage
	if err := json.Unmarshal(msg.Body, &request); err != nil {
		slog.Error("Failed to unmarshal legacy request", "error", err)
		cm.sendErrorResponse(ctx, "", fmt.Errorf("invalid message format: %w", err))
		msg.Nack(false, false)
		return
	}

	guildID := request.GetGuildID()
	slog.Info("Processing legacy event", "correlation_id", request.CorrelationID, "guild_id", guildID, "event_type", request.EventType)

	// Validate guild_id
	if guildID == "" {
		slog.Error("Missing guild_id in legacy request")
		cm.sendErrorResponse(ctx, request.CorrelationID, fmt.Errorf("guild_id is required"))
		msg.Nack(false, false)
		return
	}

	// Get handler for event type
	handler, ok := cm.handlers[request.EventType]
	if !ok {
		slog.Error("Unsupported event type in legacy queue", "event_type", request.EventType)
		cm.sendErrorResponse(ctx, request.CorrelationID, fmt.Errorf("unsupported event type: %s", request.EventType))
		msg.Nack(false, false)
		return
	}

	// Prepare payload based on event type
	payload := request.Payload
	if isApplicationEvent(request.EventType) {
		payload = map[string]interface{}{
			"event_type":              string(request.EventType),
			"contest_id":              request.ContestID,
			"user_id":                 request.UserID,
			"discord_user_id":         request.DiscordUserID,
			"discord_guild_id":        request.DiscordGuildID,
			"discord_text_channel_id": request.DiscordTextChannelID,
			"data":                    request.Data,
		}
	} else if isTeamEvent(request.EventType) {
		var fullPayload map[string]interface{}
		if err := json.Unmarshal(msg.Body, &fullPayload); err != nil {
			slog.Error("Failed to unmarshal team event payload", "error", err)
			cm.sendErrorResponse(ctx, request.CorrelationID, fmt.Errorf("invalid team event payload: %w", err))
			msg.Nack(false, false)
			return
		}
		payload = fullPayload
	}

	// Handle the event
	data, err := handler.Handle(ctx, cm.bot, guildID, payload)
	if err != nil {
		slog.Error("Legacy handler failed", "correlation_id", request.CorrelationID, "error", err)
		cm.sendErrorResponse(ctx, request.CorrelationID, err)
		msg.Nack(false, false)
		return
	}

	// Send success response
	cm.sendSuccessResponse(ctx, request.CorrelationID, data)
	msg.Ack(false)
	slog.Info("Legacy event processed successfully", "correlation_id", request.CorrelationID)
}

// resolveEventType extracts the event type from AMQP headers first, then falls back to JSON body.
func (cm *ConsumerManager) resolveEventType(msg amqp.Delivery) EventType {
	// Try AMQP header first
	if msg.Headers != nil {
		if et, ok := msg.Headers["event_type"]; ok {
			if s, ok := et.(string); ok && s != "" {
				return EventType(s)
			}
		}
	}

	// Fallback: parse JSON body for event_type field
	var body struct {
		EventType string `json:"event_type"`
	}
	if err := json.Unmarshal(msg.Body, &body); err == nil && body.EventType != "" {
		return EventType(body.EventType)
	}

	return ""
}

// extractGuildID extracts guild_id from a payload map, trying discord_guild_id first.
func extractGuildID(payload map[string]interface{}) string {
	if gid, ok := payload["discord_guild_id"].(string); ok && gid != "" {
		return gid
	}
	if gid, ok := payload["guild_id"].(string); ok && gid != "" {
		return gid
	}
	return ""
}

// sendSuccessResponse sends a success response via the publisher
func (cm *ConsumerManager) sendSuccessResponse(ctx context.Context, correlationID string, data map[string]interface{}) {
	if cm.publisher == nil {
		return
	}
	response := &ResponseMessage{
		CorrelationID: correlationID,
		Success:       true,
		Data:          data,
	}
	if err := cm.publisher.Publish(ctx, response); err != nil {
		slog.Error("Failed to publish success response", "correlation_id", correlationID, "error", err)
	}
}

// sendErrorResponse sends an error response via the publisher
func (cm *ConsumerManager) sendErrorResponse(ctx context.Context, correlationID string, err error) {
	if cm.publisher == nil {
		return
	}
	response := &ResponseMessage{
		CorrelationID: correlationID,
		Success:       false,
		Error:         err.Error(),
	}
	if pubErr := cm.publisher.Publish(ctx, response); pubErr != nil {
		slog.Error("Failed to publish error response", "correlation_id", correlationID, "error", pubErr)
	}
}

// isApplicationEvent checks if the event type is an application event
func isApplicationEvent(eventType EventType) bool {
	switch eventType {
	case EventApplicationRequested, EventApplicationAccepted, EventApplicationRejected, EventApplicationCancelled:
		return true
	default:
		return false
	}
}

// isTeamEvent checks if the event type is a team event
func isTeamEvent(eventType EventType) bool {
	switch eventType {
	case EventTeamInviteSent, EventTeamInviteAccepted, EventTeamInviteRejected,
		EventTeamMemberJoined, EventTeamMemberLeft, EventTeamMemberKicked,
		EventTeamLeadershipTransferred, EventTeamFinalized, EventTeamDeleted:
		return true
	default:
		return false
	}
}

// Close closes all channels managed by the ConsumerManager
func (cm *ConsumerManager) Close() error {
	var lastErr error
	for _, ch := range cm.channels {
		if ch != nil {
			if err := ch.Close(); err != nil {
				lastErr = err
			}
		}
	}
	return lastErr
}
