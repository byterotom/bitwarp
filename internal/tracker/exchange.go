package tracker

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

// function to setup exchange
func ExchangeInit() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("error connecting rabbit server: %v", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatalf("error establishing channel: %v", err)
	}

	err = ch.ExchangeDeclare(
		"bitwarp", // name
		"fanout",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("error declaring exchange: %v", err)
	}

	log.Printf("exchange [bitwarp] started")
}

// function to stop exchange
func StopExchange() {
	ch.Close()
	conn.Close()
}

// function to publish resource request
func PublishRequest(data string) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := ch.PublishWithContext(ctx,
		"bitwarp", // exchange
		"",        // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		},
	)
	if err != nil {
		log.Fatalf("error publishing request to nodes: %v", err)
	}
}
