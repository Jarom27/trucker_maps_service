package connection

import (
	"fmt"
	"net/http"
	"trucker_maps_service/adapters"
)

func StartServer(wsAdapter adapters.AdapterType, host string, port string) {
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := adapters.NewWebSocket(wsAdapter, w, r)
		if err != nil {
			fmt.Println("Failed to upgrade WebSocket:", err)
			return
		}
		defer conn.Close()

		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading WebSocket message:", err)
				break
			}
			fmt.Printf("Message: %s", msg)
			//service.HandleMessage(string(msg))
		}
	})

	fmt.Println("Server started on port:", port)
	address := fmt.Sprintf("%s:%s", host, port)
	http.ListenAndServe(address, nil)
}
