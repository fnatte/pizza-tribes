package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type ClientInitHandler = func(ctx context.Context, c *Client) error
type ClientMessageHandler = func(ctx context.Context, m []byte, c *Client)
type AuthFunc = func(r *http.Request) error

type WsHandler interface {
	HandleInit(ctx context.Context, c *Client) error
	HandleMessage(ctx context.Context, m []byte, c *Client)
}

type WsEndpoint struct {
	authFunc AuthFunc
	upgrader websocket.Upgrader
	hub      *Hub
	handler  WsHandler
}

func NewEndpoint(authFunc AuthFunc, hub *Hub, handler WsHandler) *WsEndpoint {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			if r.Header.Get("Origin") == "http://localhost:3000" {
				return true
			}

			log.Warn().Msg(r.Header.Get("Origin") + "is not allowed")

			return false
		},
	}

	return &WsEndpoint{
		hub:      hub,
		authFunc: authFunc,
		upgrader: upgrader,
		handler:  handler,
	}
}

func (e *WsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ws, err := e.upgrader.Upgrade(w, r, nil)
	if err != nil {
		if _, ok := err.(websocket.HandshakeError); !ok {
			log.Error().Err(err).Msg("")
		}
		log.Warn().Err(err).Msg("")
		return
	}

	// The browser does not allow reading the HTTP response status code on
	// Web Sockets because then it could be used to probe non-ws endpoints.
	// We make the authentication check after upgrading so that we can
	// close it with a specific error code that the web app can read.
	err = e.authFunc(r)
	if err != nil {
		log.Info().Err(err).Msg("Closing websocket: unauthorized")
		ws.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(4010, "Unauthorized"),
			time.Time{},
		)
		ws.Close()
		return
	}

	userId, ok := r.Context().Value("userId").(string)
	if !ok {
		log.Warn().Msg("Failed to get account id")
		ws.Close()
		return
	}

	client := &Client{
		hub:    e.hub,
		ws:     ws,
		userId: userId,
		send:   make(chan []byte, 512),
	}
	client.hub.register <- client

	err = e.handler.HandleInit(r.Context(), client)
	if err != nil {
		log.Error().Err(err).
			Str("userId", userId).
			Msg("Failed to websocket client initialzier failed")
		ws.Close()
		return
	}

	go client.writer()
	client.reader(r.Context(), e.handler.HandleMessage)
}
