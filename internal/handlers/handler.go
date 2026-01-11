package handlers

import (
	"context"

	"github.com/gamers-bot/internal/bot"
)

// Handler defines the interface for event handlers
type Handler interface {
	// Handle processes an event and returns the result data or an error
	Handle(ctx context.Context, bot *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error)
}
