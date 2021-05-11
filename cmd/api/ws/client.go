package ws

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Maximum size in bytes for a message from the client.
	readLimit = 512

	// Time allowed to write the file to the client.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the client.
	pongWait = 60 * time.Second

	// Send pings to client with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	hub *Hub

	// The websocket connection.
	ws *websocket.Conn

	userId string

	// Buffered channel of outbound messages.
	send chan []byte
}

// Read pump loop. Reads messages from the web socket connection and invokes
// the specified handler whenever a message is received.
func (c *Client) reader(ctx context.Context, handler ClientMessageHandler) {
	defer c.ws.Close()
	c.ws.SetReadLimit(readLimit)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		messageType, message, err := c.ws.ReadMessage()
		if err != nil {
			break
		}
		if messageType != websocket.TextMessage {
			log.Warn().Msg("Received non text message")
			continue
		}

		handler(ctx, message, c)
	}
}

// Write pump loop. Writes bytes from the send chan to the actual
// web socket connection.
func (c *Client) writer() {
	pingTicker := time.NewTicker(pingPeriod)
	defer func() {
		pingTicker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))

			if !ok {
				// The hub closed the channel.
				c.ws.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.ws.WriteMessage(websocket.TextMessage, message)

		case <-pingTicker.C:
			c.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.ws.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// Send the specified bytes to the web socket connection by placing the bytes
// on the send channel. The writer() pump will pick up these bytes and send them
// on the web socket connection.
func (c *Client) Send(b []byte) {
	c.send <- b
}

// Get the user id of this client
func (c *Client) UserId() string {
	return c.userId
}
