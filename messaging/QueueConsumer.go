package messaging

import (
	"fmt"
	"time"

	"github.com/streadway/amqp"
)

type QueueConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	user    string
	pass    string
	host    string
	port    string
}

func NewQueueConsumer(queue string, user string, pass string, host string, port string) (*QueueConsumer, error) {
	connection_string := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(connection_string)
		if err == nil {
			break
		}
		fmt.Printf("Attempt %d: Failed to connect to RabbitMQ: %v\n", i, err)
		time.Sleep(time.Duration(1<<i) * time.Second) // Exponential backoff
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
	}

	channel, err = conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %w", err)
	}

	_, err = channel.QueueDeclare(
		queue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare a queue: %w", err)
	}

	fmt.Println("Successfully connected to RabbitMQ")
	return &QueueConsumer{
		conn:    conn,
		channel: channel,
		queue:   queue,
		user:    user,
		pass:    pass,
		host:    host,
		port:    port,
	}, nil
}

func (qc *QueueConsumer) Receive() (<-chan amqp.Delivery, error) {
	return qc.consume()
}

func (qc *QueueConsumer) consume() (<-chan amqp.Delivery, error) {
	msgs, err := qc.channel.Consume(
		qc.queue,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return nil, fmt.Errorf("failed to start consuming messages: %w", err)
	}

	return msgs, nil
}

func (q *QueueConsumer) Close() {
	if q.channel != nil {
		q.channel.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
	fmt.Println("RabbitMQ connection closed")
}
