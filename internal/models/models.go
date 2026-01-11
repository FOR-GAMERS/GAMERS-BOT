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
