# GAMERS-BOT

FOR GAMERS, WITH GAMERS

Discord Bot for managing voice channels and messaging with event-driven RabbitMQ architecture.

## Features

- Event-driven microservice architecture with RabbitMQ
- **Multi-guild support** - One bot instance can manage multiple Discord servers
- **Resilient RabbitMQ connection** - Bot continues operating even when RabbitMQ is down
- **Slash commands** - `/author` to show bot author, `/status` to check RabbitMQ status
- Message sending to Discord channels
- **Contest invitations** - Send contest notifications with user mentions
- Voice channel member management
- Channel information queries
- Asynchronous request/response pattern
- Automatic RabbitMQ reconnection
- Docker support for easy deployment
- Designed for horizontal scaling
- Structured JSON logging

## Architecture

This bot is designed as part of a microservice architecture (MSA). It consumes events from RabbitMQ queues and publishes results back to response queues.

```
┌─────────┐    publish event    ┌──────────────┐    consume    ┌──────────────┐    Discord API    ┌─────────┐
│   WAS   │ ─────────────────► │  RabbitMQ    │ ──────────► │ Discord Bot  │ ────────────────► │ Discord │
│         │                     │  (Request Q) │              │  (Consumer)  │                   │ Server  │
└─────────┘                     └──────────────┘              └──────────────┘                   └─────────┘
     ▲                                                                │
     │                          ┌──────────────┐                     │
     └───────────────────────── │  RabbitMQ    │ ◄───────────────────┘
          consume result         │ (Response Q) │     publish result
                                 └──────────────┘
```

## Project Structure

```
GAMERS-BOT/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── bot/
│   │   └── discord.go          # Discord bot client and handlers
│   ├── rabbitmq/
│   │   ├── consumer.go         # RabbitMQ consumer logic
│   │   ├── publisher.go        # RabbitMQ publisher for responses
│   │   └── messages.go         # Message type definitions
│   ├── handlers/
│   │   ├── handler.go          # Handler interface
│   │   ├── message.go          # Message send handler
│   │   ├── voice.go            # Voice channel operations handler
│   │   └── channel.go          # Channel info query handler
│   ├── config/
│   │   └── config.go           # Configuration management
│   └── models/
│       └── models.go           # Data models
├── docker/
│   └── Dockerfile
├── .env.example
├── go.mod
└── README.md
```

## Quick Start

```bash
# 1. Clone the repository
git clone https://github.com/yourusername/GAMERS-BOT.git
cd GAMERS-BOT

# 2. Initial setup (creates .env, installs dependencies)
make setup

# 3. Edit .env with your Discord bot token
nano .env  # or vim .env

# 4. Start RabbitMQ
make rabbitmq-start

# 5. Run the bot
make run
```

That's it! Your bot is now running. Use `make help` to see all available commands.

## Prerequisites

- Go 1.23 or higher
- Docker and Docker Compose (for containerized deployment)
- Make (for using Makefile commands)
- Discord Bot Token

## Getting Started

### 1. Create a Discord Bot

1. Go to [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Go to "Bot" section and click "Add Bot"
4. Copy the bot token
5. Enable the following Privileged Gateway Intents:
   - SERVER MEMBERS INTENT
   - MESSAGE CONTENT INTENT
6. Go to "OAuth2" > "URL Generator"
7. Select scopes: `bot` and `applications.commands`
8. Select bot permissions:
   - View Channels
   - Send Messages
   - Move Members
   - Use Slash Commands
9. Copy the generated URL and invite the bot to your server(s)

### 2. Setup RabbitMQ

Using Docker:

```bash
docker run -d --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3-management
```

Access RabbitMQ Management UI at http://localhost:15672 (default credentials: guest/guest)

### 3. Configure Environment Variables

```bash
cp .env.example .env
```

Edit `.env` file:

```env
DISCORD_TOKEN=your_discord_bot_token_here
RABBITMQ_URL=amqp://guest:guest@localhost:5672/
RABBITMQ_REQUEST_QUEUE=discord.commands
RABBITMQ_RESPONSE_QUEUE=discord.responses
RABBITMQ_PREFETCH_COUNT=1
```

**Note:** The bot supports multiple Discord servers (guilds) dynamically. Each RabbitMQ message should include the `guild_id` field to specify which server to target.

### 4. Run Locally

#### Using Makefile (Recommended)

```bash
# Initial setup (creates .env, downloads dependencies)
make setup

# Start RabbitMQ
make rabbitmq-start

# Run the bot
make run

# Or run in development mode with hot reload
make dev
```

#### Manual Setup

```bash
# Install dependencies
go mod download

# Build
go build -o bin/bot ./cmd

# Run
./bin/bot
```

Or run directly:

```bash
go run ./cmd/main.go
```

#### Available Make Commands

Run `make help` to see all available commands:

```bash
make help
```

Common commands:
- `make build` - Build the application
- `make run` - Run the application
- `make test` - Run tests
- `make docker-build` - Build Docker image
- `make rabbitmq-start` - Start RabbitMQ
- `make clean` - Clean build artifacts

## Slash Commands

The bot provides the following slash commands:

### /author

Shows the bot author information.

**Usage:** `/author`

**Response:** `Author: **SONU**`

### /status

Checks the current RabbitMQ connection status.

**Usage:** `/status`

**Response:**
- `RabbitMQ Status: 🟢 Connected` - When RabbitMQ is connected
- `RabbitMQ Status: 🔴 Disconnected` - When RabbitMQ is not connected

## Supported Events

The bot supports the following event types. All events require a `guild_id` field to specify which Discord server to target.

### SEND_MESSAGE

Send a message to a Discord channel.

**Request:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
  "guild_id": "999999999999999999",
  "event_type": "SEND_MESSAGE",
  "payload": {
    "channel_id": "123456789012345678",
    "content": "Hello, Discord!"
  }
}
```

**Response:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": true,
  "data": {
    "message_id": "987654321098765432",
    "timestamp": "2025-01-08T12:34:56Z"
  }
}
```

### MOVE_MEMBERS

Move members between voice channels.

**Request:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440001",
  "guild_id": "999999999999999999",
  "event_type": "MOVE_MEMBERS",
  "payload": {
    "from_channel_id": "111111111111111111",
    "to_channel_id": "222222222222222222",
    "user_ids": []
  }
}
```

Note: Empty `user_ids` array moves all users in the source channel.

**Response:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440001",
  "success": true,
  "data": {
    "moved_count": 5,
    "failed_users": []
  }
}
```

### GET_VOICE_CHANNELS

Get all voice channels in the guild.

**Request:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440002",
  "guild_id": "999999999999999999",
  "event_type": "GET_VOICE_CHANNELS",
  "payload": {}
}
```

**Response:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440002",
  "success": true,
  "data": {
    "channels": [
      {"id": "111111111111111111", "name": "General Voice"},
      {"id": "222222222222222222", "name": "Gaming"}
    ]
  }
}
```

### GET_TEXT_CHANNELS

Get all text channels in the guild.

**Request:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440003",
  "guild_id": "999999999999999999",
  "event_type": "GET_TEXT_CHANNELS",
  "payload": {}
}
```

**Response:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440003",
  "success": true,
  "data": {
    "channels": [
      {"id": "333333333333333333", "name": "general"},
      {"id": "444444444444444444", "name": "announcements"}
    ]
  }
}
```

### SEND_CONTEST_INVITATION

Send a contest invitation with user mentions.

**Request:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440004",
  "guild_id": "999999999999999999",
  "event_type": "SEND_CONTEST_INVITATION",
  "payload": {
    "channel_id": "333333333333333333",
    "user_ids": ["111111111111111111", "222222222222222222"],
    "contest_name": "Algorithm Challenge 2025",
    "message": "Get ready for the most exciting coding challenge!"
  }
}
```

Note: The `message` field is optional. If not provided, a default message will be used.

**Response:**
```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440004",
  "success": true,
  "data": {
    "message_id": "987654321098765432",
    "notified_users": ["111111111111111111", "222222222222222222"],
    "timestamp": "2025-01-08T12:34:56Z"
  }
}
```

### Error Response

When an error occurs:

```json
{
  "correlation_id": "550e8400-e29b-41d4-a716-446655440000",
  "success": false,
  "error": "channel not found: 123456789012345678"
}
```

## Publishing Events

### Using RabbitMQ Management UI

1. Go to http://localhost:15672
2. Login with credentials (default: guest/guest)
3. Navigate to "Queues" tab
4. Click on `discord.commands` queue
5. Expand "Publish message" section
6. Set "Delivery mode" to "2 - Persistent"
7. Paste your JSON payload
8. Click "Publish message"

### Using Python

```python
import pika
import json
import uuid

# Connect to RabbitMQ
connection = pika.BlockingConnection(
    pika.ConnectionParameters('localhost')
)
channel = connection.channel()

# Declare queues
channel.queue_declare(queue='discord.commands', durable=True)
channel.queue_declare(queue='discord.responses', durable=True)

# Publish event
event = {
    "correlation_id": str(uuid.uuid4()),
    "guild_id": "999999999999999999",
    "event_type": "SEND_MESSAGE",
    "payload": {
        "channel_id": "123456789012345678",
        "content": "Hello from Python!"
    }
}

channel.basic_publish(
    exchange='',
    routing_key='discord.commands',
    body=json.dumps(event),
    properties=pika.BasicProperties(
        delivery_mode=2,  # make message persistent
    )
)

print(f"Published event with correlation_id: {event['correlation_id']}")

# Consume response
def callback(ch, method, properties, body):
    response = json.loads(body)
    if response['correlation_id'] == event['correlation_id']:
        print(f"Received response: {response}")
        ch.stop_consuming()

channel.basic_consume(
    queue='discord.responses',
    on_message_callback=callback,
    auto_ack=True
)

print("Waiting for response...")
channel.start_consuming()

connection.close()
```

### Using Go

```go
package main

import (
    "context"
    "encoding/json"
    "log"
    "time"

    "github.com/google/uuid"
    amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    ch, err := conn.Channel()
    if err != nil {
        log.Fatal(err)
    }
    defer ch.Close()

    // Publish event
    correlationID := uuid.New().String()
    event := map[string]interface{}{
        "correlation_id": correlationID,
        "guild_id":       "999999999999999999",
        "event_type":     "SEND_MESSAGE",
        "payload": map[string]interface{}{
            "channel_id": "123456789012345678",
            "content":    "Hello from Go!",
        },
    }

    body, _ := json.Marshal(event)
    err = ch.PublishWithContext(
        context.Background(),
        "",                // exchange
        "discord.commands", // routing key
        false,             // mandatory
        false,             // immediate
        amqp.Publishing{
            DeliveryMode: amqp.Persistent,
            ContentType:  "application/json",
            Body:         body,
        },
    )
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Published event with correlation_id: %s", correlationID)
}
```

## Development

### Quick Start with Makefile

```bash
# Setup project
make setup

# Start RabbitMQ
make rabbitmq-start

# Run in development mode (with hot reload)
make dev

# In another terminal, check RabbitMQ status
make rabbitmq-ui
```

### Build

```bash
# Using Make
make build

# Or manually
go build -o bin/bot ./cmd
```

### Run Tests

```bash
# Using Make
make test

# With coverage
make test-coverage

# Or manually
go test ./...
```

### Code Quality

```bash
# Format code
make fmt

# Run all checks (format, vet, lint)
make check

# Or individually
go fmt ./...
golangci-lint run
```

### Install Development Tools

```bash
# Install golangci-lint and air (hot reload)
make install-tools
```

## Docker

### Build Image

```bash
docker build -f docker/Dockerfile -t gamers-discord-bot .
```

### Run Container

```bash
docker run -d \
  --name gamers-bot \
  -e DISCORD_TOKEN=your_token \
  -e RABBITMQ_URL=amqp://guest:guest@rabbitmq:5672/ \
  --link rabbitmq \
  gamers-discord-bot
```

### Docker Compose

#### Using Makefile (Recommended)

```bash
# Start all services (RabbitMQ + Bot)
make compose-up

# View logs
make compose-logs

# Stop all services
make compose-down

# Rebuild and start
make compose-build
```

#### Manual Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild
docker-compose up -d --build
```

The included `docker-compose.yml` provides:
- RabbitMQ with management UI
- Discord Bot with automatic reconnection
- Health checks
- Persistent volumes for RabbitMQ data
- Automatic restart on failure

## Troubleshooting

### Bot is not responding

1. Check if the bot is online in your Discord server
2. Verify the `DISCORD_TOKEN` is correct
3. Ensure the bot has proper permissions
4. Check the logs: `docker-compose logs -f discord-bot`

### RabbitMQ connection failed

**Important:** The bot will continue to operate even if RabbitMQ is not available. It will automatically attempt to reconnect every 10 seconds.

1. Verify RabbitMQ is running: `docker ps | grep rabbitmq`
2. Check the `RABBITMQ_URL` is correct
3. Ensure the bot can reach RabbitMQ (network connectivity)
4. Use `/status` slash command in Discord to check connection status
5. Check bot logs for reconnection attempts

### Events not being processed

1. Check if events are published to the correct queue (`discord.commands`)
2. Verify the event format matches the examples
3. Check bot logs for errors
4. Verify the correlation_id is a valid string

### Cannot move members

1. Ensure the bot has "Move Members" permission
2. Verify that users are actually in the source voice channel
3. Check if the destination voice channel exists
4. Bot cannot move users who have higher roles than the bot

## Multi-Guild Support

The bot is designed to handle multiple Discord servers (guilds) dynamically. To use this feature:

1. **Invite the bot to multiple servers** using the OAuth2 URL from the Discord Developer Portal
2. **Include guild_id in each RabbitMQ message** to specify which server to target
3. **Get Guild ID**: Enable Developer Mode in Discord (Settings > Advanced > Developer Mode), then right-click your server and select "Copy ID"

The bot will automatically handle requests for any guild it has been invited to, without requiring configuration changes or restarts.

## RabbitMQ Resilience

The bot is designed to be resilient to RabbitMQ connection issues:

### Automatic Reconnection

- The bot starts even if RabbitMQ is unavailable
- Automatically attempts to reconnect every 10 seconds
- No manual intervention required

### Connection Monitoring

- Use `/status` slash command to check current RabbitMQ connection status
- Logs all connection attempts and failures
- Status updates are logged in real-time

### Graceful Degradation

- Discord slash commands (`/author`, `/status`) work independently of RabbitMQ
- Bot remains responsive to Discord events
- WAS integration is automatically restored when RabbitMQ comes back online

### Best Practices

1. Monitor bot logs for connection status
2. Use the `/status` command to verify connectivity
3. Ensure RabbitMQ has proper health checks in production
4. Consider setting up alerts for prolonged disconnections

## Future Enhancements

- Web dashboard for monitoring
- Webhook support for event notifications
- Metrics and logging integration (Prometheus, Grafana)
- Guild whitelist/blacklist for security
- Dead letter queue for failed events
- Contest management with leaderboards
- Scheduled contest reminders
- User participation tracking

## License

This project is licensed under the MIT License.

## Support

For issues and questions, please open an issue on GitHub.
