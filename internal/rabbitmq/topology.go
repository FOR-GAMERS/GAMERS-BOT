package rabbitmq

// QueueBinding defines a queue name and its routing key bindings
type QueueBinding struct {
	QueueName   string
	RoutingKeys []string
}

// DefaultQueueBindings returns the unified notification queue binding for the gamers.events exchange.
// All notification events are consumed from a single dedicated queue: discord.bot.notifications.
func DefaultQueueBindings() []QueueBinding {
	return []QueueBinding{
		{
			QueueName: "discord.bot.notifications",
			RoutingKeys: []string{
				"contest.#",
				"game.team.#",
				"game.game.#",
			},
		},
	}
}
