package bot

import (
	"fmt"
	"log/slog"

	"github.com/bwmarrin/discordgo"
	"github.com/gamers-bot/internal/models"
)

// DiscordBot wraps the Discord session and provides helper methods
type DiscordBot struct {
	Session              *discordgo.Session
	ready                chan struct{}
	rabbitMQConnected    bool
	statusNotificationCh chan string
}

// New creates a new Discord bot instance
func New(token string) (*DiscordBot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("failed to create Discord session: %w", err)
	}

	bot := &DiscordBot{
		Session:              session,
		ready:                make(chan struct{}),
		rabbitMQConnected:    false,
		statusNotificationCh: make(chan string, 10),
	}

	// Register event handlers
	session.AddHandler(bot.onReady)
	session.AddHandler(bot.onInteractionCreate)

	// Set intents
	session.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildVoiceStates

	return bot, nil
}

// Connect establishes connection to Discord
func (b *DiscordBot) Connect() error {
	if err := b.Session.Open(); err != nil {
		return fmt.Errorf("failed to open Discord session: %w", err)
	}
	return nil
}

// WaitUntilReady blocks until the bot is ready
func (b *DiscordBot) WaitUntilReady() {
	<-b.ready
}

// Close closes the Discord session
func (b *DiscordBot) Close() error {
	return b.Session.Close()
}

// onReady is called when the bot is ready
func (b *DiscordBot) onReady(s *discordgo.Session, event *discordgo.Ready) {
	slog.Info("Discord bot is ready", "user", event.User.Username)

	// Register slash commands
	if err := b.RegisterCommands(); err != nil {
		slog.Error("Failed to register commands", "error", err)
	}

	close(b.ready)
}

// onInteractionCreate handles slash command interactions
func (b *DiscordBot) onInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	switch i.ApplicationCommandData().Name {
	case "author":
		b.handleAuthorCommand(s, i)
	case "status":
		b.handleStatusCommand(s, i)
	case "damepo":
		b.handleDamepoCommand(s, i)
	}
}

// handleAuthorCommand responds with "SONU"
func (b *DiscordBot) handleAuthorCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "**SONU**っていう人から生まれました",
		},
	})
	if err != nil {
		slog.Error("Failed to respond to author command", "error", err)
	}
}

func (b *DiscordBot) handleDamepoCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "キモオジだめぽ😩",
		},
	})
	if err != nil {
		slog.Error("Failed to respond to author command", "error", err)
	}
}

// handleStatusCommand responds with RabbitMQ connection status
func (b *DiscordBot) handleStatusCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	status := "🔴 Disconnected"
	if b.rabbitMQConnected {
		status = "🟢 Connected"
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("RabbitMQ Status: %s", status),
		},
	})
	if err != nil {
		slog.Error("Failed to respond to status command", "error", err)
	}
}

// RegisterCommands registers slash commands with Discord
func (b *DiscordBot) RegisterCommands() error {
	commands := []*discordgo.ApplicationCommand{
		{
			Name:        "author",
			Description: "Show the bot author",
		},
		{
			Name:        "status",
			Description: "Check RabbitMQ connection status",
		},
	}

	for _, cmd := range commands {
		_, err := b.Session.ApplicationCommandCreate(b.Session.State.User.ID, "", cmd)
		if err != nil {
			return fmt.Errorf("failed to create command %s: %w", cmd.Name, err)
		}
		slog.Info("Registered command", "command", cmd.Name)
	}

	return nil
}

// NotifyRabbitMQStatus updates the RabbitMQ connection status
func (b *DiscordBot) NotifyRabbitMQStatus(connected bool, err error) {
	b.rabbitMQConnected = connected

	var message string
	if connected {
		message = "✅ RabbitMQ connection established"
		slog.Info(message)
	} else {
		message = fmt.Sprintf("⚠️ RabbitMQ connection failed: %v", err)
		slog.Warn(message, "error", err)
	}

	// Send notification to channel (non-blocking)
	select {
	case b.statusNotificationCh <- message:
	default:
		// Channel is full, skip notification
	}
}

// SendMessage sends a message to a Discord channel
func (b *DiscordBot) SendMessage(channelID, content string) (*models.SendMessageResult, error) {
	message, err := b.Session.ChannelMessageSend(channelID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	return &models.SendMessageResult{
		MessageID: message.ID,
		Timestamp: message.Timestamp.Format("2006-01-02T15:04:05Z"),
	}, nil
}

// MoveMembers moves members between voice channels
func (b *DiscordBot) MoveMembers(guildID, fromChannelID, toChannelID string, userIDs []string) (*models.MoveMembersResult, error) {
	// Validate that channels are different
	if fromChannelID == toChannelID {
		return nil, fmt.Errorf("source and destination channels must be different")
	}

	// Get members in the source voice channel
	guild, err := b.Session.State.Guild(guildID)
	if err != nil {
		// If not in state, fetch from API
		guild, err = b.Session.Guild(guildID)
		if err != nil {
			return nil, fmt.Errorf("failed to get guild: %w", err)
		}
	}

	// Collect members to move
	var membersToMove []string
	if len(userIDs) > 0 {
		// Move specific users
		membersToMove = userIDs
	} else {
		// Move all users in the source channel
		for _, vs := range guild.VoiceStates {
			if vs.ChannelID == fromChannelID {
				membersToMove = append(membersToMove, vs.UserID)
			}
		}
	}

	if len(membersToMove) == 0 {
		return &models.MoveMembersResult{
			MovedCount:  0,
			FailedUsers: []string{},
		}, nil
	}

	// Move members
	var failedUsers []string
	movedCount := 0

	for _, userID := range membersToMove {
		err := b.Session.GuildMemberMove(guildID, userID, &toChannelID)
		if err != nil {
			slog.Warn("Failed to move user", "user_id", userID, "error", err)
			failedUsers = append(failedUsers, userID)
		} else {
			movedCount++
		}
	}

	return &models.MoveMembersResult{
		MovedCount:  movedCount,
		FailedUsers: failedUsers,
	}, nil
}

// GetVoiceChannels retrieves all voice channels in the guild
func (b *DiscordBot) GetVoiceChannels(guildID string) (*models.GetChannelsResult, error) {
	channels, err := b.Session.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild channels: %w", err)
	}

	var voiceChannels []models.ChannelInfo
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildVoice {
			voiceChannels = append(voiceChannels, models.ChannelInfo{
				ID:   channel.ID,
				Name: channel.Name,
			})
		}
	}

	return &models.GetChannelsResult{
		Channels: voiceChannels,
	}, nil
}

// GetTextChannels retrieves all text channels in the guild
func (b *DiscordBot) GetTextChannels(guildID string) (*models.GetChannelsResult, error) {
	channels, err := b.Session.GuildChannels(guildID)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild channels: %w", err)
	}

	var textChannels []models.ChannelInfo
	for _, channel := range channels {
		if channel.Type == discordgo.ChannelTypeGuildText {
			textChannels = append(textChannels, models.ChannelInfo{
				ID:   channel.ID,
				Name: channel.Name,
			})
		}
	}

	return &models.GetChannelsResult{
		Channels: textChannels,
	}, nil
}

// SendContestInvitation sends a contest invitation to specified users
func (b *DiscordBot) SendContestInvitation(channelID string, userIDs []string, contestName, customMessage string) (*models.ContestInvitationResult, error) {
	// Build mention string
	var mentions string
	for _, userID := range userIDs {
		mentions += fmt.Sprintf("<@%s> ", userID)
	}

	// Build the message
	var content string
	if customMessage != "" {
		content = fmt.Sprintf("🎮 **Contest Invitation: %s**\n\n%s\n\n%s", contestName, customMessage, mentions)
	} else {
		content = fmt.Sprintf("🎮 **Contest Invitation: %s**\n\n%sYou have been invited to participate in this contest!", contestName, mentions)
	}

	// Send the message
	message, err := b.Session.ChannelMessageSend(channelID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to send contest invitation: %w", err)
	}

	return &models.ContestInvitationResult{
		MessageID:     message.ID,
		NotifiedUsers: userIDs,
		Timestamp:     message.Timestamp.Format("2006-01-02T15:04:05Z"),
	}, nil
}
