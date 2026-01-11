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

					// Initialize publisher
					publisher, err := rabbitmq.NewPublisher(conn, cfg.RabbitMQResponseQueue)
					if err != nil {
						slog.Error("Failed to create publisher", "error", err)
						conn.Close()
						time.Sleep(10 * time.Second)
						continue
					}

					slog.Info("Publisher initialized", "queue", cfg.RabbitMQResponseQueue)

					// Initialize consumer
					consumer, err := rabbitmq.NewConsumer(
						conn,
						cfg.RabbitMQRequestQueue,
						cfg.RabbitMQPrefetchCount,
						discordBot,
						publisher,
					)
					if err != nil {
						slog.Error("Failed to create consumer", "error", err)
						publisher.Close()
						conn.Close()
						time.Sleep(10 * time.Second)
						continue
					}

					// Register handlers
					consumer.RegisterHandler(rabbitmq.EventSendMessage, handlers.NewMessageHandler())
					consumer.RegisterHandler(rabbitmq.EventMoveMembers, handlers.NewVoiceHandler())
					consumer.RegisterHandler(rabbitmq.EventGetVoiceChannels, handlers.NewVoiceChannelHandler())
					consumer.RegisterHandler(rabbitmq.EventGetTextChannels, handlers.NewTextChannelHandler())
					consumer.RegisterHandler(rabbitmq.EventSendContestInvitation, handlers.NewContestInvitationHandler())

					slog.Info("Handlers registered")

					// Start consumer (blocking)
					slog.Info("Starting consumer", "queue", cfg.RabbitMQRequestQueue)
					if err := consumer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
						slog.Error("Consumer error", "error", err)
					}

					// Cleanup
					consumer.Close()
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
