package adapters

import (
	"github.com/gorilla/websocket"
)

// GorillaAdapter implementa WebSocketConnection usando Gorilla WebSocket.
type GorillaAdapter struct {
	Conn *websocket.Conn
}

func (g *GorillaAdapter) WriteMessage(messageType int, data []byte) error {
	return g.Conn.WriteMessage(messageType, data)
}

func (g *GorillaAdapter) ReadMessage() (messageType int, data []byte, err error) {
	return g.Conn.ReadMessage()
}

func (g *GorillaAdapter) Close() error {
	return g.Conn.Close()
}
