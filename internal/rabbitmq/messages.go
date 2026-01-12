package rabbitmq

// EventType represents the type of event being processed
type EventType string

const (
	// EventSendMessage sends a message to a Discord channel
	EventSendMessage EventType = "SEND_MESSAGE"

	// EventMoveMembers moves members between voice channels
	EventMoveMembers EventType = "MOVE_MEMBERS"

	// EventGetVoiceChannels retrieves all voice channels in the guild
	EventGetVoiceChannels EventType = "GET_VOICE_CHANNELS"

	// EventGetTextChannels retrieves all text channels in the guild
	EventGetTextChannels EventType = "GET_TEXT_CHANNELS"

	// EventSendContestInvitation sends a contest invitation to users
	EventSendContestInvitation EventType = "SEND_CONTEST_INVITATION"

	// EventApplicationRequested notifies that a user has requested to join a contest
	EventApplicationRequested EventType = "application.requested"

	// EventApplicationAccepted notifies a user that their contest application was accepted
	EventApplicationAccepted EventType = "application.accepted"

	// EventApplicationRejected notifies a user that their contest application was rejected
	EventApplicationRejected EventType = "application.rejected"
)

// RequestMessage represents an incoming event from the request queue
// Supports both legacy format (guild_id) and WAS format (discord_guild_id)
type RequestMessage struct {
	CorrelationID        string                 `json:"correlation_id"`
	GuildID              string                 `json:"guild_id"`
	DiscordGuildID       string                 `json:"discord_guild_id"`
	DiscordTextChannelID string                 `json:"discord_text_channel_id"`
	DiscordUserID        string                 `json:"discord_user_id"`
	ContestID            int64                  `json:"contest_id"`
	UserID               int64                  `json:"user_id"`
	EventType            EventType              `json:"event_type"`
	Payload              map[string]interface{} `json:"payload"`
	Data                 map[string]interface{} `json:"data"`
}

// GetGuildID returns the guild ID from either field
func (r *RequestMessage) GetGuildID() string {
	if r.DiscordGuildID != "" {
		return r.DiscordGuildID
	}
	return r.GuildID
}

// ResponseMessage represents an outgoing response to the response queue
type ResponseMessage struct {
	CorrelationID string                 `json:"correlation_id"`
	Success       bool                   `json:"success"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Error         string                 `json:"error,omitempty"`
}
