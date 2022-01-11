package models

import (
	"kshoplistSrv/constants"
	"time"

	"github.com/gorilla/websocket"
)

// connection is an middleman between the websocket connection and the hub.
type Connection struct {
	// The websocket connection.
	Ws *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte
}

// write writes a message with the given message type and payload.
func (c *Connection) Write(mt int, payload []byte) error {
	c.Ws.SetWriteDeadline(time.Now().Add(constants.WriteWait))
	return c.Ws.WriteMessage(mt, payload)
}
