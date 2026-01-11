package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// MessageHandler handles SEND_MESSAGE events
type MessageHandler struct{}

// NewMessageHandler creates a new MessageHandler
func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

// Handle processes a SEND_MESSAGE event
func (h *MessageHandler) Handle(ctx context.Context, bot *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// Parse payload to SendMessagePayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var sendPayload models.SendMessagePayload
	if err := json.Unmarshal(payloadBytes, &sendPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Validate payload
	if sendPayload.ChannelID == "" {
		return nil, fmt.Errorf("channel_id is required")
	}
	if sendPayload.Content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// Send message
	result, err := bot.SendMessage(sendPayload.ChannelID, sendPayload.Content)
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
