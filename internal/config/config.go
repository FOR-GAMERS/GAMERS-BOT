package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DiscordToken string

	RabbitMQURL           string
	RabbitMQRequestQueue  string
	RabbitMQResponseQueue string
	RabbitMQPrefetchCount int
}

func Load() (*Config, error) {
	_ = godotenv.Load("env/.env")

	// Try RABBITMQ_URL first, then build from individual env vars
	// RabbitMQ is optional - if not configured, bot runs without it
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = buildRabbitMQURL()
	}

	config := &Config{
		DiscordToken:          os.Getenv("DISCORD_TOKEN"),
		RabbitMQURL:           rabbitMQURL,
		RabbitMQRequestQueue:  getEnvOrDefault("RABBITMQ_REQUEST_QUEUE", "discord.commands"),
		RabbitMQResponseQueue: getEnvOrDefault("RABBITMQ_RESPONSE_QUEUE", "discord.responses"),
		RabbitMQPrefetchCount: getEnvAsIntOrDefault("RABBITMQ_PREFETCH_COUNT", 1),
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// Validate checks if required configuration values are set
func (c *Config) Validate() error {
	if c.DiscordToken == "" {
		return fmt.Errorf("DISCORD_TOKEN is required")
	}
	// RabbitMQ is optional - only validate prefetch count if RabbitMQ is enabled
	if c.RabbitMQEnabled() && c.RabbitMQPrefetchCount < 1 {
		return fmt.Errorf("RABBITMQ_PREFETCH_COUNT must be at least 1")
	}
	return nil
}

// RabbitMQEnabled returns true if RabbitMQ is configured
func (c *Config) RabbitMQEnabled() bool {
	return c.RabbitMQURL != ""
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsIntOrDefault returns the value of an environment variable as an int or a default value
func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// buildRabbitMQURL constructs a RabbitMQ URL from individual environment variables
// Format: amqp://user:password@host:port/vhost
func buildRabbitMQURL() string {
	host := getEnvOrDefault("RABBITMQ_HOST", "")
	port := getEnvOrDefault("RABBITMQ_PORT", "")
	user := getEnvOrDefault("RABBITMQ_USER", "")
	password := getEnvOrDefault("RABBITMQ_PASSWORD", "")
	vhost := getEnvOrDefault("RABBITMQ_VHOST", "")

	return fmt.Sprintf("amqp://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		vhost,
	)
}
