package rabbitmq

// QueueBinding defines a queue name and its routing key bindings
type QueueBinding struct {
	QueueName   string
	RoutingKeys []string
}

// DefaultQueueBindings returns the 4 new queue bindings for the gamers.events exchange
func DefaultQueueBindings() []QueueBinding {
	return []QueueBinding{
		{
			QueueName: "bot.contest.notifications",
			RoutingKeys: []string{
				"contest.#",
			},
		},
		{
			QueueName: "bot.team.notifications",
			RoutingKeys: []string{
				"game.team.#",
			},
		},
		{
			QueueName: "bot.game.notifications",
			RoutingKeys: []string{
				"game.scheduled",
				"game.activated",
				"game.match.*",
				"game.finished",
			},
		},
		{
			QueueName: "bot.contest.teams.ready",
			RoutingKeys: []string{
				"game.contest.teams.ready",
			},
		},
	}
}
