package adapters

import (
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
)

// Tipo de adaptador
type AdapterType string

const (
	Gorilla AdapterType = "gorilla"
)

// NewWebSocket crea una nueva conexión WebSocket basada en el adaptador seleccionado.
func NewWebSocket(adapterType AdapterType, w http.ResponseWriter, r *http.Request) (WebSocketConnection, error) {
	switch adapterType {
	case Gorilla:
		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Permite todas las conexiones (en producción, restringe esto adecuadamente)
			},
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return nil, err
		}
		return &GorillaAdapter{Conn: conn}, nil

	default:
		return nil, ErrUnsupportedAdapter
	}
}

// ErrUnsupportedAdapter se lanza cuando el adaptador no es compatible.
var ErrUnsupportedAdapter = errors.New("unsupported WebSocket adapter")
