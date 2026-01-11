package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// ContestInvitationHandler handles SEND_CONTEST_INVITATION events
type ContestInvitationHandler struct{}

// NewContestInvitationHandler creates a new ContestInvitationHandler
func NewContestInvitationHandler() *ContestInvitationHandler {
	return &ContestInvitationHandler{}
}

// Handle processes a SEND_CONTEST_INVITATION event
func (h *ContestInvitationHandler) Handle(ctx context.Context, bot *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	// Parse payload to ContestInvitationPayload
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var invitePayload models.ContestInvitationPayload
	if err := json.Unmarshal(payloadBytes, &invitePayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	// Validate payload
	if invitePayload.ChannelID == "" {
		return nil, fmt.Errorf("channel_id is required")
	}
	if len(invitePayload.UserIDs) == 0 {
		return nil, fmt.Errorf("user_ids cannot be empty")
	}
	if invitePayload.ContestName == "" {
		return nil, fmt.Errorf("contest_name is required")
	}

	// Send contest invitation
	result, err := bot.SendContestInvitation(
		invitePayload.ChannelID,
		invitePayload.UserIDs,
		invitePayload.ContestName,
		invitePayload.Message,
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
