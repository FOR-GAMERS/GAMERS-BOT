package models

// ChannelInfo represents basic information about a Discord channel
type ChannelInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// SendMessagePayload contains parameters for sending a message
type SendMessagePayload struct {
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
}

// SendMessageResult contains the result of sending a message
type SendMessageResult struct {
	MessageID string `json:"message_id"`
	Timestamp string `json:"timestamp"`
}

// MoveMembersPayload contains parameters for moving members between channels
type MoveMembersPayload struct {
	FromChannelID string   `json:"from_channel_id"`
	ToChannelID   string   `json:"to_channel_id"`
	UserIDs       []string `json:"user_ids"` // Empty = move all users
}

// MoveMembersResult contains the result of moving members
type MoveMembersResult struct {
	MovedCount  int      `json:"moved_count"`
	FailedUsers []string `json:"failed_users"`
}

// GetChannelsResult contains a list of channels
type GetChannelsResult struct {
	Channels []ChannelInfo `json:"channels"`
}

// ContestInvitationPayload contains parameters for sending a contest invitation
type ContestInvitationPayload struct {
	ChannelID   string   `json:"channel_id"`   // Text channel to send the invitation
	UserIDs     []string `json:"user_ids"`     // User IDs to mention
	ContestName string   `json:"contest_name"` // Name of the contest
	Message     string   `json:"message"`      // Optional custom message
}

// ContestInvitationResult contains the result of sending a contest invitation
type ContestInvitationResult struct {
	MessageID     string   `json:"message_id"`
	NotifiedUsers []string `json:"notified_users"`
	Timestamp     string   `json:"timestamp"`
}

// ContestApplicationEventPayload represents the event payload from Web Application Server
type ContestApplicationEventPayload struct {
	EventType            string                 `json:"event_type"`
	ContestID            int64                  `json:"contest_id"`
	UserID               int64                  `json:"user_id"`
	DiscordUserID        string                 `json:"discord_user_id"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	Data                 map[string]interface{} `json:"data"`
}

// ApplicationNotificationResult contains the result of sending an application notification
type ApplicationNotificationResult struct {
	MessageID string `json:"message_id"`
	UserID    string `json:"user_id"`
	Timestamp string `json:"timestamp"`
}

// TeamInviteEventPayload represents the payload for team invite events
// Used for: team.invite.sent, team.invite.accepted, team.invite.rejected
type TeamInviteEventPayload struct {
	EventID              string                 `json:"event_id"`
	EventType            string                 `json:"event_type"`
	Timestamp            string                 `json:"timestamp"`
	GameID               int64                  `json:"game_id"`
	ContestID            int64                  `json:"contest_id"`
	InviterUserID        int64                  `json:"inviter_user_id"`
	InviterDiscordID     string                 `json:"inviter_discord_id"`
	InviterUsername      string                 `json:"inviter_username"`
	InviteeUserID        int64                  `json:"invitee_user_id"`
	InviteeDiscordID     string                 `json:"invitee_discord_id"`
	InviteeUsername      string                 `json:"invitee_username"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	TeamName             string                 `json:"team_name"`
	Data                 map[string]interface{} `json:"data"`
}

// TeamMemberEventPayload represents the payload for team member events
// Used for: team.member.joined, team.member.left, team.member.kicked
type TeamMemberEventPayload struct {
	EventID              string                 `json:"event_id"`
	EventType            string                 `json:"event_type"`
	Timestamp            string                 `json:"timestamp"`
	GameID               int64                  `json:"game_id"`
	ContestID            int64                  `json:"contest_id"`
	UserID               int64                  `json:"user_id"`
	DiscordUserID        string                 `json:"discord_user_id"`
	Username             string                 `json:"username"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	CurrentMemberCount   int                    `json:"current_member_count"`
	MaxMembers           int                    `json:"max_members"`
	Data                 map[string]interface{} `json:"data"`
}

// TeamFinalizedEventPayload represents the payload for team finalized/deleted/leadership events
// Used for: team.finalized, team.deleted, team.leadership.transferred
type TeamFinalizedEventPayload struct {
	EventID              string                 `json:"event_id"`
	EventType            string                 `json:"event_type"`
	Timestamp            string                 `json:"timestamp"`
	GameID               int64                  `json:"game_id"`
	ContestID            int64                  `json:"contest_id"`
	LeaderUserID         int64                  `json:"leader_user_id"`
	LeaderDiscordID      string                 `json:"leader_discord_id"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	MemberCount          int                    `json:"member_count"`
	MemberUserIDs        []int64                `json:"member_user_ids"`
	Data                 map[string]interface{} `json:"data"`
}

// TeamNotificationResult contains the result of sending a team notification
type TeamNotificationResult struct {
	MessageID string `json:"message_id"`
	Timestamp string `json:"timestamp"`
}

// BaseEvent contains common fields embedded in all event payloads from gamers.events
type BaseEvent struct {
	EventID   string `json:"event_id"`
	EventType string `json:"event_type"`
	Timestamp string `json:"timestamp"`
}

// GameEventPayload represents the payload for game lifecycle events
// Used for: game.scheduled, game.activated, game.match.*, game.finished
type GameEventPayload struct {
	BaseEvent
	GameID               int64                  `json:"game_id"`
	ContestID            int64                  `json:"contest_id"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	Data                 map[string]interface{} `json:"data"`
}

// ContestTeamsReadyPayload represents the payload for game.contest.teams.ready events
type ContestTeamsReadyPayload struct {
	BaseEvent
	GameID               int64                  `json:"game_id"`
	ContestID            int64                  `json:"contest_id"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	TeamCount            int                    `json:"team_count"`
	Data                 map[string]interface{} `json:"data"`
}

// ContestCreatedEventPayload represents the payload for contest.created events
type ContestCreatedEventPayload struct {
	BaseEvent
	ContestID            int64                  `json:"contest_id"`
	ContestTitle         string                 `json:"contest_title"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	Data                 map[string]interface{} `json:"data"`
}
