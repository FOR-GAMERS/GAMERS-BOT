package handlers

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/models"
)

// ==================== Team Invite Handlers ====================

// TeamInviteSentHandler handles team.invite.sent events
type TeamInviteSentHandler struct{}

// NewTeamInviteSentHandler creates a new TeamInviteSentHandler
func NewTeamInviteSentHandler() *TeamInviteSentHandler {
	return &TeamInviteSentHandler{}
}

// Handle processes a team.invite.sent event - sends DM to invitee
func (h *TeamInviteSentHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamInvitePayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.InviteeDiscordID == "" {
		return nil, fmt.Errorf("invitee_discord_id is required")
	}
	if eventPayload.TeamName == "" {
		return nil, fmt.Errorf("team_name is required")
	}

	// Build DM content
	content := fmt.Sprintf(
		"**[チーム招待]**\n\n"+
			"**%s**さんから **%s** チームに招待されました。\n"+
			"招待を確認して、参加するかどうかを決めてください。",
		eventPayload.InviterUsername, eventPayload.TeamName,
	)

	// Send DM to invitee
	result, err := b.SendDirectMessage(eventPayload.InviteeDiscordID, content)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// TeamInviteAcceptedHandler handles team.invite.accepted events
type TeamInviteAcceptedHandler struct{}

// NewTeamInviteAcceptedHandler creates a new TeamInviteAcceptedHandler
func NewTeamInviteAcceptedHandler() *TeamInviteAcceptedHandler {
	return &TeamInviteAcceptedHandler{}
}

// Handle processes a team.invite.accepted event - sends message to team channel
func (h *TeamInviteAcceptedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamInvitePayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordTextChannelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}
	if eventPayload.InviteeDiscordID == "" {
		return nil, fmt.Errorf("invitee_discord_id is required")
	}

	// Send notification to team channel
	result, err := b.SendTeamInviteNotification(
		eventPayload.DiscordTextChannelID,
		eventPayload.InviterDiscordID,
		eventPayload.InviterUsername,
		eventPayload.InviteeDiscordID,
		eventPayload.InviteeUsername,
		eventPayload.TeamName,
		bot.TeamInviteAccepted,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// TeamInviteRejectedHandler handles team.invite.rejected events
type TeamInviteRejectedHandler struct{}

// NewTeamInviteRejectedHandler creates a new TeamInviteRejectedHandler
func NewTeamInviteRejectedHandler() *TeamInviteRejectedHandler {
	return &TeamInviteRejectedHandler{}
}

// Handle processes a team.invite.rejected event - sends DM to inviter (team leader)
func (h *TeamInviteRejectedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamInvitePayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.InviterDiscordID == "" {
		return nil, fmt.Errorf("inviter_discord_id is required")
	}

	// Build DM content for team leader
	content := fmt.Sprintf(
		"**[招待拒否]**\n\n"+
			"**%s**さんが **%s** チームへの招待を拒否しました。",
		eventPayload.InviteeUsername, eventPayload.TeamName,
	)

	// Send DM to inviter (team leader)
	result, err := b.SendDirectMessage(eventPayload.InviterDiscordID, content)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// ==================== Team Member Handlers ====================

// TeamMemberJoinedHandler handles team.member.joined events
type TeamMemberJoinedHandler struct{}

// NewTeamMemberJoinedHandler creates a new TeamMemberJoinedHandler
func NewTeamMemberJoinedHandler() *TeamMemberJoinedHandler {
	return &TeamMemberJoinedHandler{}
}

// Handle processes a team.member.joined event - sends welcome message to team channel
func (h *TeamMemberJoinedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamMemberPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordTextChannelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}
	if eventPayload.DiscordUserID == "" {
		return nil, fmt.Errorf("discord_user_id is required")
	}

	// Send notification to team channel
	result, err := b.SendTeamMemberNotification(
		eventPayload.DiscordTextChannelID,
		eventPayload.DiscordUserID,
		eventPayload.Username,
		eventPayload.CurrentMemberCount,
		eventPayload.MaxMembers,
		bot.TeamMemberJoined,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// TeamMemberLeftHandler handles team.member.left events
type TeamMemberLeftHandler struct{}

// NewTeamMemberLeftHandler creates a new TeamMemberLeftHandler
func NewTeamMemberLeftHandler() *TeamMemberLeftHandler {
	return &TeamMemberLeftHandler{}
}

// Handle processes a team.member.left event - sends notification to team channel
func (h *TeamMemberLeftHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamMemberPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordTextChannelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}

	// Send notification to team channel
	result, err := b.SendTeamMemberNotification(
		eventPayload.DiscordTextChannelID,
		eventPayload.DiscordUserID,
		eventPayload.Username,
		eventPayload.CurrentMemberCount,
		eventPayload.MaxMembers,
		bot.TeamMemberLeft,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// TeamMemberKickedHandler handles team.member.kicked events
type TeamMemberKickedHandler struct{}

// NewTeamMemberKickedHandler creates a new TeamMemberKickedHandler
func NewTeamMemberKickedHandler() *TeamMemberKickedHandler {
	return &TeamMemberKickedHandler{}
}

// Handle processes a team.member.kicked event - sends DM to kicked user
func (h *TeamMemberKickedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamMemberPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordUserID == "" {
		return nil, fmt.Errorf("discord_user_id is required")
	}

	// Build DM content for kicked user
	content := "**[チーム強制退出]**\n\n" +
		"チームから退出されました。\n" +
		"詳しい内容はチームリーダーにお問い合わせください。"

	// Send DM to kicked user
	result, err := b.SendDirectMessage(eventPayload.DiscordUserID, content)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// ==================== Team Status Handlers ====================

// TeamLeadershipTransferredHandler handles team.leadership.transferred events
type TeamLeadershipTransferredHandler struct{}

// NewTeamLeadershipTransferredHandler creates a new TeamLeadershipTransferredHandler
func NewTeamLeadershipTransferredHandler() *TeamLeadershipTransferredHandler {
	return &TeamLeadershipTransferredHandler{}
}

// Handle processes a team.leadership.transferred event - sends notification to team channel
func (h *TeamLeadershipTransferredHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamFinalizedPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordTextChannelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}
	if eventPayload.LeaderDiscordID == "" {
		return nil, fmt.Errorf("leader_discord_id is required")
	}

	// Send notification to team channel
	result, err := b.SendTeamStatusNotification(
		eventPayload.DiscordTextChannelID,
		eventPayload.LeaderDiscordID,
		eventPayload.MemberCount,
		bot.TeamLeadershipTransferred,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// TeamFinalizedHandler handles team.finalized events
type TeamFinalizedHandler struct{}

// NewTeamFinalizedHandler creates a new TeamFinalizedHandler
func NewTeamFinalizedHandler() *TeamFinalizedHandler {
	return &TeamFinalizedHandler{}
}

// Handle processes a team.finalized event - sends notification to team channel
func (h *TeamFinalizedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamFinalizedPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordTextChannelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}
	if eventPayload.LeaderDiscordID == "" {
		return nil, fmt.Errorf("leader_discord_id is required")
	}

	// Send notification to team channel
	result, err := b.SendTeamStatusNotification(
		eventPayload.DiscordTextChannelID,
		eventPayload.LeaderDiscordID,
		eventPayload.MemberCount,
		bot.TeamFinalized,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// TeamDeletedHandler handles team.deleted events
type TeamDeletedHandler struct{}

// NewTeamDeletedHandler creates a new TeamDeletedHandler
func NewTeamDeletedHandler() *TeamDeletedHandler {
	return &TeamDeletedHandler{}
}

// Handle processes a team.deleted event - sends notification to team channel
func (h *TeamDeletedHandler) Handle(ctx context.Context, b *bot.DiscordBot, guildID string, payload map[string]interface{}) (map[string]interface{}, error) {
	eventPayload, err := parseTeamFinalizedPayload(payload)
	if err != nil {
		return nil, err
	}

	// Validate required fields
	if eventPayload.DiscordTextChannelID == "" {
		return nil, fmt.Errorf("discord_text_channel_id is required")
	}

	// Send notification to team channel
	result, err := b.SendTeamStatusNotification(
		eventPayload.DiscordTextChannelID,
		eventPayload.LeaderDiscordID,
		eventPayload.MemberCount,
		bot.TeamDeleted,
	)
	if err != nil {
		return nil, err
	}

	return marshalResult(result)
}

// ==================== Helper Functions ====================

// parseTeamInvitePayload parses the payload into TeamInviteEventPayload
func parseTeamInvitePayload(payload map[string]interface{}) (*models.TeamInviteEventPayload, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var eventPayload models.TeamInviteEventPayload
	if err := json.Unmarshal(payloadBytes, &eventPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &eventPayload, nil
}

// parseTeamMemberPayload parses the payload into TeamMemberEventPayload
func parseTeamMemberPayload(payload map[string]interface{}) (*models.TeamMemberEventPayload, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var eventPayload models.TeamMemberEventPayload
	if err := json.Unmarshal(payloadBytes, &eventPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &eventPayload, nil
}

// parseTeamFinalizedPayload parses the payload into TeamFinalizedEventPayload
func parseTeamFinalizedPayload(payload map[string]interface{}) (*models.TeamFinalizedEventPayload, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %w", err)
	}

	var eventPayload models.TeamFinalizedEventPayload
	if err := json.Unmarshal(payloadBytes, &eventPayload); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payload: %w", err)
	}

	return &eventPayload, nil
}

// marshalResult converts a result struct to map[string]interface{}
func marshalResult(result interface{}) (map[string]interface{}, error) {
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
