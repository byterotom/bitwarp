package node

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel
var q amqp.Queue

// function to declare queue
func QueueInit() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("error connecting rabbit server: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("error establishing channel: %v", err)
	}

	q, err = ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("error declaring queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name,    // queue name
		"",        // routing key
		"bitwarp", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("error binding queue: %v", err)
	}
}

// function to stop queue
func StopQueue() {
	ch.Close()
	conn.Close()
}

// function to recieve message
func ConsumeMessage() {
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			log.Printf(" [x] %s", d.Body)
		}
	}()

	if err != nil {
		log.Fatalf("error publishing message: %v", err)
	}
}
