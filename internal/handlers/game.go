package handlers

import (
	"context"
	"log/slog"

	"github.com/gamers-bot/internal/bot"
)

// GameScheduledHandler handles game.scheduled events
type GameScheduledHandler struct{}

func NewGameScheduledHandler() *GameScheduledHandler {
	return &GameScheduledHandler{}
}

func (h *GameScheduledHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("GameScheduledHandler invoked", "guild_id", guildID)
	return nil, nil
}

// GameActivatedHandler handles game.activated events
type GameActivatedHandler struct{}

func NewGameActivatedHandler() *GameActivatedHandler {
	return &GameActivatedHandler{}
}

func (h *GameActivatedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("GameActivatedHandler invoked", "guild_id", guildID)
	return nil, nil
}

// GameMatchDetectingHandler handles game.match.detecting events
type GameMatchDetectingHandler struct{}

func NewGameMatchDetectingHandler() *GameMatchDetectingHandler {
	return &GameMatchDetectingHandler{}
}

func (h *GameMatchDetectingHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("GameMatchDetectingHandler invoked", "guild_id", guildID)
	return nil, nil
}

// GameMatchDetectedHandler handles game.match.detected events
type GameMatchDetectedHandler struct{}

func NewGameMatchDetectedHandler() *GameMatchDetectedHandler {
	return &GameMatchDetectedHandler{}
}

func (h *GameMatchDetectedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("GameMatchDetectedHandler invoked", "guild_id", guildID)
	return nil, nil
}

// GameMatchFailedHandler handles game.match.failed events
type GameMatchFailedHandler struct{}

func NewGameMatchFailedHandler() *GameMatchFailedHandler {
	return &GameMatchFailedHandler{}
}

func (h *GameMatchFailedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("GameMatchFailedHandler invoked", "guild_id", guildID)
	return nil, nil
}

// GameFinishedHandler handles game.finished events
type GameFinishedHandler struct{}

func NewGameFinishedHandler() *GameFinishedHandler {
	return &GameFinishedHandler{}
}

func (h *GameFinishedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("GameFinishedHandler invoked", "guild_id", guildID)
	return nil, nil
}

// ContestTeamsReadyHandler handles game.contest.teams.ready events
type ContestTeamsReadyHandler struct{}

func NewContestTeamsReadyHandler() *ContestTeamsReadyHandler {
	return &ContestTeamsReadyHandler{}
}

func (h *ContestTeamsReadyHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// TODO: Discord 알림 전송 로직
	slog.Info("ContestTeamsReadyHandler invoked", "guild_id", guildID)
	return nil, nil
}
