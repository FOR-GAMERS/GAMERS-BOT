package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// ContestCreatedHandler handles contest.created events
type ContestCreatedHandler struct{}

func NewContestCreatedHandler() *ContestCreatedHandler {
	return &ContestCreatedHandler{}
}

func (h *ContestCreatedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	var event models.ContestCreatedEventPayload
	if err := unmarshalPayload(payload, &event); err != nil {
		return nil, err
	}

	if event.DiscordTextChannelID == "" {
		slog.Warn("contest.created: discord_text_channel_id is empty, skipping", "contest_id", event.ContestID)
		return nil, nil
	}

	content := fmt.Sprintf(
		"**[大会作成]**\n\n"+
			"新しい大会 **%s** が作成されました！\n"+
			"参加申請をお待ちしております。",
		event.ContestTitle,
	)

	result, err := b.SendMessage(event.DiscordTextChannelID, content)
	if err != nil {
		return nil, fmt.Errorf("contest.created: failed to send message: %w", err)
	}

	return marshalResult(result)
}

// ContestInvitationHandler handles SEND_CONTEST_INVITATION events
type ContestInvitationHandler struct{}

// NewContestInvitationHandler creates a new ContestInvitationHandler
func NewContestInvitationHandler() *ContestInvitationHandler {
	return &ContestInvitationHandler{}
}

// Handle processes a SEND_CONTEST_INVITATION event
func (h *ContestInvitationHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var invitePayload models.ContestInvitationPayload
	if err := json.Unmarshal(payloadBytes, &invitePayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	if invitePayload.ChannelID == "" {
		return nil, fmt.Errorf("channel_id is required")
	}
	if len(invitePayload.UserIDs) == 0 {
		return nil, fmt.Errorf("user_ids cannot be empty")
	}
	if invitePayload.ContestName == "" {
		return nil, fmt.Errorf("contest_name is required")
	}

	result, err := b.SendContestInvitation(
		invitePayload.ChannelID,
		invitePayload.UserIDs,
		invitePayload.ContestName,
		invitePayload.Message,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// unmarshalPayload is a helper to marshal+unmarshal a map into a typed struct.
func unmarshalPayload(payload map[string]interface{}, dest interface{}) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}
	if err := json.Unmarshal(b, dest); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	return nil
}
