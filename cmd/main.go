package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gamers-bot/internal/bot"
	"github.com/gamers-bot/internal/config"
	"github.com/gamers-bot/internal/handlers"
	"github.com/gamers-bot/internal/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Setup pretty logging with charmbracelet/log
	logger := log.NewWithOptions(os.Stdout, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		Level:           log.InfoLevel,
	})
	slog.SetDefault(slog.New(logger))

	slog.Info("Starting GAMERS Discord Bot")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	slog.Info("Configuration loaded successfully")

	// Initialize Discord bot
	discordBot, err := bot.New(cfg.DiscordToken)
	if err != nil {
		slog.Error("Failed to create Discord bot", "error", err)
		os.Exit(1)
	}

	// Connect to Discord
	if err := discordBot.Connect(); err != nil {
		slog.Error("Failed to connect to Discord", "error", err)
		os.Exit(1)
	}
	defer discordBot.Close()

	slog.Info("Connected to Discord, waiting for bot to be ready...")
	discordBot.WaitUntilReady()
	slog.Info("Discord bot is ready")

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Try to connect to RabbitMQ (non-blocking) - only if enabled
	if cfg.RabbitMQEnabled() {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					slog.Info("Attempting to connect to RabbitMQ", "url", cfg.RabbitMQURL)
					conn, err := amqp.Dial(cfg.RabbitMQURL)
					if err != nil {
						slog.Warn("Failed to connect to RabbitMQ, will retry in 10 seconds", "error", err)
						discordBot.NotifyRabbitMQStatus(false, err)
						time.Sleep(10 * time.Second)
						continue
					}

					slog.Info("Connected to RabbitMQ")
					discordBot.NotifyRabbitMQStatus(true, nil)

					// Initialize publisher (for legacy queue responses)
					publisher, err := rabbitmq.NewPublisher(conn, cfg.RabbitMQResponseQueue)
					if err != nil {
						slog.Error("Failed to create publisher", "error", err)
						conn.Close()
						time.Sleep(10 * time.Second)
						continue
					}

					slog.Info("Publisher initialized", "queue", cfg.RabbitMQResponseQueue)

					// Initialize ConsumerManager
					manager := rabbitmq.NewConsumerManager(
						conn,
						cfg.RabbitMQExchange,
						cfg.RabbitMQPrefetchCount,
						discordBot,
						publisher,
					)

					// Setup new queue topology (gamers.events exchange + 4 queues)
					if err := manager.SetupTopology(); err != nil {
						slog.Error("Failed to setup topology", "error", err)
						publisher.Close()
						conn.Close()
						time.Sleep(10 * time.Second)
						continue
					}

					// Setup legacy queue (discord.commands bound to legacy exchanges)
					legacyBindings := []rabbitmq.LegacyBinding{
						{Exchange: cfg.RabbitMQTeamExchange, RoutingKey: cfg.RabbitMQTeamRoutingKey},
					}
					if err := manager.SetupLegacyQueue(cfg.RabbitMQRequestQueue, legacyBindings); err != nil {
						slog.Error("Failed to setup legacy queue", "error", err)
						publisher.Close()
						conn.Close()
						time.Sleep(10 * time.Second)
						continue
					}

					// Register legacy handlers (request/response pattern)
					manager.RegisterHandler(rabbitmq.EventSendMessage, handlers.NewMessageHandler())
					manager.RegisterHandler(rabbitmq.EventMoveMembers, handlers.NewVoiceHandler())
					manager.RegisterHandler(rabbitmq.EventGetVoiceChannels, handlers.NewVoiceChannelHandler())
					manager.RegisterHandler(rabbitmq.EventGetTextChannels, handlers.NewTextChannelHandler())
					manager.RegisterHandler(rabbitmq.EventSendContestInvitation, handlers.NewContestInvitationHandler())

					// Register application event handlers
					manager.RegisterHandler(rabbitmq.EventApplicationRequested, handlers.NewApplicationRequestedHandler())
					manager.RegisterHandler(rabbitmq.EventApplicationAccepted, handlers.NewApplicationAcceptedHandler())
					manager.RegisterHandler(rabbitmq.EventApplicationRejected, handlers.NewApplicationRejectedHandler())
					manager.RegisterHandler(rabbitmq.EventApplicationCancelled, handlers.NewApplicationCancelledHandler())
					manager.RegisterHandler(rabbitmq.EventMemberWithdrawn, handlers.NewMemberWithdrawnHandler())

					// Register team event handlers
					manager.RegisterHandler(rabbitmq.EventTeamInviteSent, handlers.NewTeamInviteSentHandler())
					manager.RegisterHandler(rabbitmq.EventTeamInviteAccepted, handlers.NewTeamInviteAcceptedHandler())
					manager.RegisterHandler(rabbitmq.EventTeamInviteRejected, handlers.NewTeamInviteRejectedHandler())
					manager.RegisterHandler(rabbitmq.EventTeamMemberJoined, handlers.NewTeamMemberJoinedHandler())
					manager.RegisterHandler(rabbitmq.EventTeamMemberLeft, handlers.NewTeamMemberLeftHandler())
					manager.RegisterHandler(rabbitmq.EventTeamMemberKicked, handlers.NewTeamMemberKickedHandler())
					manager.RegisterHandler(rabbitmq.EventTeamLeadershipTransferred, handlers.NewTeamLeadershipTransferredHandler())
					manager.RegisterHandler(rabbitmq.EventTeamFinalized, handlers.NewTeamFinalizedHandler())
					manager.RegisterHandler(rabbitmq.EventTeamDeleted, handlers.NewTeamDeletedHandler())

					// Register contest event handlers
					manager.RegisterHandler(rabbitmq.EventContestCreated, handlers.NewContestCreatedHandler())

					// Register game event handlers
					manager.RegisterHandler(rabbitmq.EventGameScheduled, handlers.NewGameScheduledHandler())
					manager.RegisterHandler(rabbitmq.EventGameActivated, handlers.NewGameActivatedHandler())
					manager.RegisterHandler(rabbitmq.EventGameMatchDetecting, handlers.NewGameMatchDetectingHandler())
					manager.RegisterHandler(rabbitmq.EventGameMatchDetected, handlers.NewGameMatchDetectedHandler())
					manager.RegisterHandler(rabbitmq.EventGameMatchFailed, handlers.NewGameMatchFailedHandler())
					manager.RegisterHandler(rabbitmq.EventGameFinished, handlers.NewGameFinishedHandler())

					// Register contest teams ready handler
					manager.RegisterHandler(rabbitmq.EventContestTeamsReady, handlers.NewContestTeamsReadyHandler())

					slog.Info("All handlers registered")

					// Start all consumers (blocking)
					slog.Info("Starting ConsumerManager",
						"exchange", cfg.RabbitMQExchange,
						"legacy_queue", cfg.RabbitMQRequestQueue,
					)
					if err := manager.Start(ctx, cfg.RabbitMQRequestQueue); err != nil && !errors.Is(err, context.Canceled) {
						slog.Error("ConsumerManager error", "error", err)
					}

					// Cleanup
					manager.Close()
					publisher.Close()
					conn.Close()

					// If context is cancelled, exit the goroutine
					if ctx.Err() != nil {
						return
					}

					// Otherwise, reconnect after delay
					slog.Info("RabbitMQ connection lost, reconnecting in 10 seconds...")
					discordBot.NotifyRabbitMQStatus(false, errors.New("connection lost"))
					time.Sleep(10 * time.Second)
				}
			}
		}()
	} else {
		slog.Info("RabbitMQ not configured, running Discord bot only")
	}

	// Wait for shutdown signal
	<-sigChan
	slog.Info("Received shutdown signal, gracefully shutting down...")
	cancel()

	// Give some time for graceful shutdown
	time.Sleep(2 * time.Second)

	slog.Info("GAMERS Discord Bot stopped")
}
