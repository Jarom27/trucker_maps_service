package main

import (
	"fmt"
	"trucker_maps_service/adapters"
	"trucker_maps_service/config"
	"trucker_maps_service/connection"
	"trucker_maps_service/domain"
	"trucker_maps_service/messaging"
)

func main() {
	config := config.LoadConfig()
	gpsManager := domain.NewGPSManager()
	dataConsumer, err := messaging.NewQueueConsumer(
		config.RabbitQueue,
		config.RabbitUser,
		config.RabbitPass,
		config.RabbitHost,
		config.RabbitPort,
		10,
	)
	if err != nil {
		fmt.Printf("there was an error when initializing consumer %s", err)
	}
	go func() {
		dataConsumer.Start(gpsManager, 10)
	}()
	connection.StartServer(gpsManager, adapters.Gorilla, config.Host, config.Port)
}
