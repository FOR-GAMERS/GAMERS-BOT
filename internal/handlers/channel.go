package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gamers-bot/internal/bot"
)

// VoiceChannelHandler handles GET_VOICE_CHANNELS events
type VoiceChannelHandler struct{}

// NewVoiceChannelHandler creates a handler for GET_VOICE_CHANNELS events
func NewVoiceChannelHandler() *VoiceChannelHandler {
	return &VoiceChannelHandler{}
}

// Handle processes a GET_VOICE_CHANNELS event
func (h *VoiceChannelHandler) Handle(ctx context.Context, bot *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	result, err := bot.GetVoiceChannels(guildID)
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

// TextChannelHandler handles GET_TEXT_CHANNELS events
type TextChannelHandler struct{}

// NewTextChannelHandler creates a handler for GET_TEXT_CHANNELS events
func NewTextChannelHandler() *TextChannelHandler {
	return &TextChannelHandler{}
}

// Handle processes a GET_TEXT_CHANNELS event
func (h *TextChannelHandler) Handle(ctx context.Context, bot *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	result, err := bot.GetTextChannels(guildID)
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
