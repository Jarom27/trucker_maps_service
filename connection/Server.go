package connection

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"trucker_maps_service/adapters"
	"trucker_maps_service/domain"

	"github.com/gorilla/websocket"
)

func StartServer(serviceManager *domain.GPSManager, wsAdapter adapters.AdapterType, host string, port string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {

		conn, err := adapters.NewWebSocket(wsAdapter, w, r)
		if err != nil {
			fmt.Println("Failed to upgrade WebSocket:", err)
			return
		}
		defer conn.Close()

		// Canal para indicar si el cliente sigue conectado
		done := make(chan bool)

		// Listener para recibir mensajes desde el cliente (si es necesario)
		go func() {
			for {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					fmt.Println("Error reading WebSocket message:", err)
					done <- true // Finalizar la conexión
					break
				}
				fmt.Printf("Received Message: %s\n", msg)
			}
		}()

		// Ticker para enviar datos periódicamente a los clientes
		ticker := time.NewTicker(2 * time.Second) // Actualización cada 2 segundos
		defer ticker.Stop()

		for {
			select {
			case <-done:
				fmt.Println("Client disconnected")
				return
			case <-ticker.C:
				// Obtener los datos GPS actuales
				gpsData := serviceManager.GetAllGPSData()
				if len(gpsData) == 0 {
					continue
				}

				// Serializar los datos a JSON
				data, err := json.Marshal(gpsData)
				if err != nil {
					fmt.Println("Failed to serialize GPS data:", err)
					continue
				}

				// Enviar los datos por WebSocket
				err = conn.WriteMessage(websocket.TextMessage, data)
				if err != nil {
					fmt.Println("Failed to send GPS data:", err)
					done <- true // Finalizar la conexión si hay error
					break
				}
			}
		}
	})

	// Iniciar el servidor HTTP
	address := fmt.Sprintf("%s:%s", host, port)
	fmt.Println("Server started on", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
