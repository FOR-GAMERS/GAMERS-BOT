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
)

// RequestMessage represents an incoming event from the request queue
type RequestMessage struct {
	CorrelationID string                 `json:"correlation_id"`
	GuildID       string                 `json:"guild_id"`
	EventType     EventType              `json:"event_type"`
	Payload       map[string]interface{} `json:"payload"`
}

// ResponseMessage represents an outgoing response to the response queue
type ResponseMessage struct {
	CorrelationID string                 `json:"correlation_id"`
	Success       bool                   `json:"success"`
	Data          map[string]interface{} `json:"data,omitempty"`
	Error         string                 `json:"error,omitempty"`
}
