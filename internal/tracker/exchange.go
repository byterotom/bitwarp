package tracker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type SyncMessage struct {
	Sender   string
	NodeIp   string
	FileHash string
	Chunks   []uint64
}

var conn *amqp.Connection
var ch *amqp.Channel

func ExchangeInit() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Fatal("failed to connect to RabbitMQ:", err)
	}

	ch, err = conn.Channel()
	if err != nil {
		log.Fatal("failed to open a channel:", err)
	}

	err = ch.ExchangeDeclare(
		"redis_sync", // exchange name
		"fanout",     // exchange type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatal("failed to declare an exchange:", err)
	}

}

func Publish(msg *SyncMessage) {
	jsonData, err := json.Marshal(msg)
	if err != nil {
		log.Fatal("failed to marshal json:", err)
	}

	// publish the message
	err = ch.Publish(
		"redis_sync", // exchange name
		"",           // routing key (empty for fanout)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        jsonData,
		},
	)
	if err != nil {
		log.Fatal("failed to publish:", err)
	}
}

func (tr *TrackerServer) Sync() {
	queue, err := ch.QueueDeclare(
		"",    // unique name
		false, // durable
		true,  // auto-delete
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		queue.Name,   // queue name
		"",           // routing key (empty for fanout)
		"redis_sync", // exchange name
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		queue.Name, // queue name
		"",         // consumer tag
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("failed to start consuming messages: %v", err)
	}

	for msg := range msgs {
		var temp SyncMessage
		err = json.Unmarshal(msg.Body, &temp)
		if err != nil {
			fmt.Printf("failed to unmarshal message: %v\n", err)
			continue
		}

		if temp.Sender == tr.address {
			continue
		}

		go func() {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			for _, chunkNo := range temp.Chunks {

				// construct key
				key := fmt.Sprintf("%s:%d", temp.FileHash, chunkNo)
				// construct score
				score := float64(time.Now().Add(TTL).Unix())

				err := Rdb.ZAdd(ctx, key, redis.Z{Score: score, Member: temp.NodeIp}).Err()
				if err != nil {
					log.Printf("error adding holder to redis: %v", err)
				}
			}

		}()

	}
}

func Stop() {
	conn.Close()
	ch.Close()
}
