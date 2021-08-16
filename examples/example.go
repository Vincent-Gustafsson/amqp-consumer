package example

import (
	"fmt"

	"github.com/streadway/amqp"
	"github.com/vincent-gustafsson/amqp-consumer"
)

func main() {
	client := consumer.NewConsumerClient()
	fmt.Println("[*] Connecting")
	client.Connect("YOUR CONN_STRING")

	dopts := &consumer.SubscribeToQueueOpts{
		QueueOpts: &consumer.QueueConfig{
			Name:       "",
			Durable:    false,
			AutoDelete: false,
			Exclusive:  false,
			NoWait:     false,
			Args:       nil,
		},
		QueueBindOpts: &consumer.QueueBindConfig{
			RoutingKey: "",
			Exchange:   "errors",
			NoWait:     false,
			Args:       nil,
		},
		ConsumeOpts: &consumer.ConsumeConfig{
			ConsumerName: "asdasd",
			Exclusive:    false,
			NoWait:       false,
			Args:         nil,
		},
	}

	forever := make(chan bool)

	fmt.Println("[*] Subscribing to queue")
	client.SubscribeToQueue(func(d *amqp.Delivery) bool {
		fmt.Println(string(d.Body))
		return true
	}, dopts)

	<-forever
}
