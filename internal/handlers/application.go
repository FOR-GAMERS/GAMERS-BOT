package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// ApplicationRequestedHandler handles APPLICATION_REQUESTED events
type ApplicationRequestedHandler struct{}

func NewApplicationRequestedHandler() *ApplicationRequestedHandler {
	return &ApplicationRequestedHandler{}
}

func (h *ApplicationRequestedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	return handleApplicationNotification(b, payload, bot.StatusRequested)
}

// ApplicationAcceptedHandler handles APPLICATION_ACCEPTED events
type ApplicationAcceptedHandler struct{}

func NewApplicationAcceptedHandler() *ApplicationAcceptedHandler {
	return &ApplicationAcceptedHandler{}
}

func (h *ApplicationAcceptedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	return handleApplicationNotification(b, payload, bot.StatusAccepted)
}

// ApplicationRejectedHandler handles APPLICATION_REJECTED events
type ApplicationRejectedHandler struct{}

func NewApplicationRejectedHandler() *ApplicationRejectedHandler {
	return &ApplicationRejectedHandler{}
}

func (h *ApplicationRejectedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	return handleApplicationNotification(b, payload, bot.StatusRejected)
}

// ApplicationCancelledHandler handles application.cancelled events
type ApplicationCancelledHandler struct{}

func NewApplicationCancelledHandler() *ApplicationCancelledHandler {
	return &ApplicationCancelledHandler{}
}

func (h *ApplicationCancelledHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.ContestApplicationEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("application.cancelled: discord_text_channel_id is empty, skipping")
		return nil, nil
	}
	if event.DiscordUserID == "" {
		slog.Warn("application.cancelled: discord_user_id is empty, skipping")
		return nil, nil
	}

	contestTitle, _ := event.Data["contest_title"].(string)
	if contestTitle == "" {
		contestTitle = "不明な大会"
	}

	content := fmt.Sprintf(
		"**[申請取消]**\n\n"+
			"<@%s>様が **%s** 大会への参加申請を取り消しました。",
		event.DiscordUserID, contestTitle,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("application.cancelled: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// MemberWithdrawnHandler handles member.withdrawn events
type MemberWithdrawnHandler struct{}

func NewMemberWithdrawnHandler() *MemberWithdrawnHandler {
	return &MemberWithdrawnHandler{}
}

func (h *MemberWithdrawnHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.ContestApplicationEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("member.withdrawn: discord_text_channel_id is empty, skipping")
		return nil, nil
	}
	if event.DiscordUserID == "" {
		slog.Warn("member.withdrawn: discord_user_id is empty, skipping")
		return nil, nil
	}

	contestTitle, _ := event.Data["contest_title"].(string)
	if contestTitle == "" {
		contestTitle = "不明な大会"
	}

	content := fmt.Sprintf(
		"**[参加辞退]**\n\n"+
			"<@%s>様が **%s** 大会への参加を辞退しました。",
		event.DiscordUserID, contestTitle,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("member.withdrawn: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// handleApplicationNotification is a shared function for handling application notifications
func handleApplicationNotification(b *bot.DiscordBot, payload map[string]interface{}, status bot.ApplicationStatus) (map[string]interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var eventPayload models.ContestApplicationEventPayload
	if err := json.Unmarshal(payloadBytes, &eventPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	channelID := eventPayload.DiscordTextChannelID
	if channelID == "" {
		slog.Warn("application notification: discord_text_channel_id is empty, skipping",
			"status", status)
		return nil, nil
	}

	if eventPayload.DiscordUserID == "" {
		slog.Warn("application notification: discord_user_id is empty, skipping",
			"status", status)
		return nil, nil
	}

	contestTitle, _ := eventPayload.Data["contest_title"].(string)
	if contestTitle == "" {
		contestTitle = "不明な大会"
	}

	var processedByDiscordID string
	if status == bot.StatusAccepted || status == bot.StatusRejected {
		processedByDiscordID, _ = eventPayload.Data["processed_by_discord_id"].(string)
	}

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

	return marshalResult(result)
}
