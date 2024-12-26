package main

import (
	"trucker_maps_service/adapters"
	"trucker_maps_service/config"
	"trucker_maps_service/connection"
)

func main() {
	config := config.LoadConfig()

	connection.StartServer(adapters.Gorilla, config.Host, config.Port)
}
