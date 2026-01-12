package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// ApplicationRequestedHandler handles APPLICATION_REQUESTED events
type ApplicationRequestedHandler struct{}

// NewApplicationRequestedHandler creates a new ApplicationRequestedHandler
func NewApplicationRequestedHandler() *ApplicationRequestedHandler {
	return &ApplicationRequestedHandler{}
}

// Handle processes an APPLICATION_REQUESTED event
func (h *ApplicationRequestedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	return handleApplicationNotification(b, payload, bot.StatusRequested)
}

// ApplicationAcceptedHandler handles APPLICATION_ACCEPTED events
type ApplicationAcceptedHandler struct{}

// NewApplicationAcceptedHandler creates a new ApplicationAcceptedHandler
func NewApplicationAcceptedHandler() *ApplicationAcceptedHandler {
	return &ApplicationAcceptedHandler{}
}

// Handle processes an APPLICATION_ACCEPTED event
func (h *ApplicationAcceptedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	return handleApplicationNotification(b, payload, bot.StatusAccepted)
}

// ApplicationRejectedHandler handles APPLICATION_REJECTED events
type ApplicationRejectedHandler struct{}

// NewApplicationRejectedHandler creates a new ApplicationRejectedHandler
func NewApplicationRejectedHandler() *ApplicationRejectedHandler {
	return &ApplicationRejectedHandler{}
}

// Handle processes an APPLICATION_REJECTED event
func (h *ApplicationRejectedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	return handleApplicationNotification(b, payload, bot.StatusRejected)
}

// handleApplicationNotification is a shared function for handling application notifications
func handleApplicationNotification(b *bot.DiscordBot, payload map[string]interface{}, status bot.ApplicationStatus) (map[string]interface{}, error) {
	// Parse payload to ContestApplicationEventPayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var eventPayload models.ContestApplicationEventPayload
	if err := json.Unmarshal(payloadBytes, &eventPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Get channel ID from DiscordTextChannelID field
	channelID := eventPayload.DiscordTextChannelID
	if channelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}

	// Validate required fields
	if eventPayload.DiscordUserID == "" {
		return nil, fmt.Errorf("discord_user_id is required")
	}

	// Extract contest_title from Data
	contestTitle, _ := eventPayload.Data["contest_title"].(string)
	if contestTitle == "" {
		return nil, fmt.Errorf("contest_title is required in data")
	}

	// Extract processed_by_discord_id for accepted/rejected events
	var processedByDiscordID string
	if status == bot.StatusAccepted || status == bot.StatusRejected {
		processedByDiscordID, _ = eventPayload.Data["processed_by_discord_id"].(string)
	}

	// Send application notification
	result, err := b.SendApplicationNotification(
		channelID,
		eventPayload.DiscordUserID,
		contestTitle,
		status,
		processedByDiscordID,
	)
	if err != nil {
		return nil, err
	}

	// Convert result to map
	resultBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal result: %w", err)
	}

	var resultMap map[string]interface{}
	if err := json.Unmarshal(resultBytes, &resultMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return resultMap, nil
}
