package consumer

import "github.com/streadway/amqp"

type MessageHandlerFunc func(*amqp.Delivery) bool

type Consumer interface {
	Connect(string) error
	SubscribeToQueue(MessageHandlerFunc, *SubscribeToQueueOpts) error
	Close()
}

type ConsumerClient struct {
	conn *amqp.Connection
}

func NewConsumerClient() *ConsumerClient {
	return &ConsumerClient{}
}

func (cc *ConsumerClient) Connect(connString string) (err error) {
	cc.conn, err = amqp.Dial(connString)
	if err != nil {
		return
	}
	return
}

func (cc *ConsumerClient) Close() {
	if cc.conn != nil {
		cc.conn.Close()
	}
}

func (cc *ConsumerClient) SubscribeToQueue(handlerFunc MessageHandlerFunc, opts *SubscribeToQueueOpts) error {
	ch, err := cc.conn.Channel()
	if err != nil {
		return err
	}

	queue, err := declareQueueWithConfig(ch, opts.QueueOpts)
	if err != nil {
		return err
	}

	err = bindQueueWithConfig(ch, queue.Name, opts.QueueBindOpts)
	if err != nil {
		return err
	}

	msgs, err := consumeWithConfig(ch, queue.Name, opts.ConsumeOpts)
	if err != nil {
		return err
	}

	go consumeMessages(msgs, handlerFunc)
	return nil
}

func consumeMessages(deliveries <-chan amqp.Delivery, handlerFunc MessageHandlerFunc) {
	for delivery := range deliveries {
		ack := handlerFunc(&delivery)
		delivery.Ack(ack)
	}
}

func declareQueueWithConfig(ch *amqp.Channel, opts *QueueConfig) (amqp.Queue, error) {
	queue, err := ch.QueueDeclare(
		opts.Name,       // name
		opts.Durable,    // durable
		opts.AutoDelete, // autoDelete
		opts.Exclusive,  // exclusive
		opts.NoWait,     // noWait
		opts.Args,       // args
	)
	if err != nil {
		return queue, err
	}
	return queue, nil
}

func bindQueueWithConfig(ch *amqp.Channel, queueName string, opts *QueueBindConfig) error {
	err := ch.QueueBind(
		queueName,       // queue name
		opts.RoutingKey, // routing key
		opts.Exchange,   // exchange
		opts.NoWait,     // noWait
		opts.Args,       // args
	)
	if err != nil {
		return err
	}
	return nil
}

func consumeWithConfig(ch *amqp.Channel, queueName string, opts *ConsumeConfig) (<-chan amqp.Delivery, error) {
	msgs, err := ch.Consume(
		queueName,         // queue
		opts.ConsumerName, // consumer
		false,             // autoAck
		opts.Exclusive,    // exclusive
		false,             // noLocal
		opts.NoWait,       // noWait
		opts.Args,         // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
