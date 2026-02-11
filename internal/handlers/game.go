package handlers

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// GameScheduledHandler handles game.scheduled events
type GameScheduledHandler struct{}

func NewGameScheduledHandler() *GameScheduledHandler {
	return &GameScheduledHandler{}
}

func (h *GameScheduledHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.GameEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.scheduled: discord_text_channel_id is empty, skipping", "game_id", event.GameID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[試合予定]**\n\n"+
			"試合(ID: %d)がスケジュールされました。\n"+
			"開始時間に備えてお待ちください！",
		event.GameID,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.scheduled: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// GameActivatedHandler handles game.activated events
type GameActivatedHandler struct{}

func NewGameActivatedHandler() *GameActivatedHandler {
	return &GameActivatedHandler{}
}

func (h *GameActivatedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.GameEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.activated: discord_text_channel_id is empty, skipping", "game_id", event.GameID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[試合開始]**\n\n"+
			"試合(ID: %d)が開始されました！\n"+
			"頑張ってください！",
		event.GameID,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.activated: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// GameMatchDetectingHandler handles game.match.detecting events
type GameMatchDetectingHandler struct{}

func NewGameMatchDetectingHandler() *GameMatchDetectingHandler {
	return &GameMatchDetectingHandler{}
}

func (h *GameMatchDetectingHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.GameEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.match.detecting: discord_text_channel_id is empty, skipping", "game_id", event.GameID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[マッチ検出中]**\n\n"+
			"試合(ID: %d)のマッチ検出を行っています...\n"+
			"しばらくお待ちください。",
		event.GameID,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.match.detecting: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// GameMatchDetectedHandler handles game.match.detected events
type GameMatchDetectedHandler struct{}

func NewGameMatchDetectedHandler() *GameMatchDetectedHandler {
	return &GameMatchDetectedHandler{}
}

func (h *GameMatchDetectedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.GameEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.match.detected: discord_text_channel_id is empty, skipping", "game_id", event.GameID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[マッチ検出完了]**\n\n"+
			"試合(ID: %d)のマッチが正常に検出されました！",
		event.GameID,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.match.detected: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// GameMatchFailedHandler handles game.match.failed events
type GameMatchFailedHandler struct{}

func NewGameMatchFailedHandler() *GameMatchFailedHandler {
	return &GameMatchFailedHandler{}
}

func (h *GameMatchFailedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.GameEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.match.failed: discord_text_channel_id is empty, skipping", "game_id", event.GameID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[マッチ検出失敗]**\n\n"+
			"試合(ID: %d)のマッチ検出に失敗しました。\n"+
			"運営人に確認をお願いします。",
		event.GameID,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.match.failed: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// GameFinishedHandler handles game.finished events
type GameFinishedHandler struct{}

func NewGameFinishedHandler() *GameFinishedHandler {
	return &GameFinishedHandler{}
}

func (h *GameFinishedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.GameEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.finished: discord_text_channel_id is empty, skipping", "game_id", event.GameID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[試合終了]**\n\n"+
			"試合(ID: %d)が終了しました。\n"+
			"お疲れ様でした！",
		event.GameID,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.finished: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// ContestTeamsReadyHandler handles game.contest.teams.ready events
type ContestTeamsReadyHandler struct{}

func NewContestTeamsReadyHandler() *ContestTeamsReadyHandler {
	return &ContestTeamsReadyHandler{}
}

func (h *ContestTeamsReadyHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.ContestTeamsReadyPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("game.contest.teams.ready: discord_text_channel_id is empty, skipping", "contest_id", event.ContestID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[チーム準備完了]**\n\n"+
			"大会の全チーム(%d チーム)が準備完了しました！\n"+
			"まもなく試合が開始されます。",
		event.TeamCount,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("game.contest.teams.ready: failed to send message: %w", err)
	}

	return marshalResult(result)
}
