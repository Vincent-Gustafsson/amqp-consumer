package consumer

import "github.com/streadway/amqp"

type SubscribeToQueueOpts struct {
	QueueOpts     *QueueConfig
	QueueBindOpts *QueueBindConfig
	ConsumeOpts   *ConsumeConfig
}

type QueueConfig struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
}

type QueueBindConfig struct {
	RoutingKey string
	Exchange   string
	NoWait     bool
	Args       amqp.Table
}

type ConsumeConfig struct {
	ConsumerName string
	Exclusive    bool
	NoWait       bool
	Args         amqp.Table
}
