package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// VoiceHandler handles MOVE_MEMBERS events
type VoiceHandler struct{}

// NewVoiceHandler creates a new VoiceHandler
func NewVoiceHandler() *VoiceHandler {
	return &VoiceHandler{}
}

// Handle processes a MOVE_MEMBERS event
func (h *VoiceHandler) Handle(ctx context.Context, bot *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// Parse payload to MoveMembersPayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var movePayload models.MoveMembersPayload
	if err := json.Unmarshal(payloadBytes, &movePayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Validate payload
	if movePayload.FromChannelID == "" {
		return nil, fmt.Errorf("from_channel_id is required")
	}
	if movePayload.ToChannelID == "" {
		return nil, fmt.Errorf("to_channel_id is required")
	}

	// Move members
	result, err := bot.MoveMembers(guildID, movePayload.FromChannelID, movePayload.ToChannelID, movePayload.UserIDs)
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
