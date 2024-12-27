package messaging

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"trucker_maps_service/domain"
	"trucker_maps_service/models"

	"github.com/streadway/amqp"
)

// QueueConsumer maneja la conexión y consumo de mensajes de RabbitMQ.
type QueueConsumer struct {
	conn       *amqp.Connection
	channel    *amqp.Channel
	gpsManager *domain.GPSManager
	queue      string
	user       string
	pass       string
	host       string
	port       string
	jobs       chan models.GPSData
}

// NewQueueConsumer inicializa la conexión con RabbitMQ.
func NewQueueConsumer(queue string, user string, pass string, host string, port string, workerPoolSize int) (*QueueConsumer, error) {
	connectionString := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, pass, host, port)
	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	const maxRetries = 5
	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(connectionString)
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

	jobs := make(chan models.GPSData, workerPoolSize*2)
	return &QueueConsumer{
		conn:    conn,
		channel: channel,
		queue:   queue,
		user:    user,
		pass:    pass,
		host:    host,
		port:    port,
		jobs:    jobs,
	}, nil
}

// Start inicia el Worker Pool y el consumo de mensajes.
func (qc *QueueConsumer) Start(gpsManager *domain.GPSManager, workerPoolSize int) {
	qc.gpsManager = gpsManager

	// Iniciar Worker Pool
	var wg sync.WaitGroup
	for i := 0; i < workerPoolSize; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for job := range qc.jobs {
				qc.gpsManager.UpdateGPSData(job.Device_id, job.Latitude, job.Longitude)
				fmt.Printf("Worker %d processed job: %+v\n", workerID, job)
			}
		}(i)
	}

	// Iniciar consumo de mensajes
	msgs, err := qc.consume()
	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	fmt.Println("Started consuming messages from RabbitMQ...")

	// Enviar mensajes al Worker Pool
	for msg := range msgs {
		var gpsMsg models.GPSData
		if err := json.Unmarshal(msg.Body, &gpsMsg); err != nil {
			log.Println("Failed to deserialize message:", err)
			continue
		}

		qc.jobs <- gpsMsg
	}

	close(qc.jobs)
	wg.Wait()
}

// consume configura el canal para consumir mensajes.
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

// Close cierra la conexión y el canal.
func (q *QueueConsumer) Close() {
	if q.channel != nil {
		q.channel.Close()
	}
	if q.conn != nil {
		q.conn.Close()
	}
	fmt.Println("RabbitMQ connection closed")
}
