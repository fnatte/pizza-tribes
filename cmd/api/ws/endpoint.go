package ws

import (
	"context"
	"net/http"
	"time"

	"github.com/fnatte/pizza-tribes/internal/game"
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

func NewEndpoint(authFunc AuthFunc, hub *Hub, handler WsHandler, origins []string) *WsEndpoint {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			for _, origin := range origins {
				if r.Header.Get("Origin") == origin {
					return true
				}
			}

			log.Warn().Msg(r.Header.Get("Origin") + " is not allowed")

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

func closeWs(ws *websocket.Conn, code int, message string) {
	ws.WriteJSON(map[string]interface{}{
		"type": "control",
		"control": map[string]interface{}{
			"type":    "close",
			"code":    code,
			"message": message,
		},
	})

	err := ws.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(code, message),
		time.Time{},
	)
	if err != nil && err != websocket.ErrCloseSent {
		ws.Close()
	}
}

func (e *WsEndpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Debug().Msg("WS Request Started")

	ws, err := e.upgrader.Upgrade(w, r, http.Header{
		"Sec-WebSocket-Protocol": []string{"pizzatribes"},
	})
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
		closeWs(ws, 4010, "unauthorized")
		return
	}

	// Read user id from request context
	userId, ok := game.GetUserIdFromContext(r.Context())
	if !ok {
		log.Warn().Msg("Closing websocket: failed to get account id")
		closeWs(ws, 5001, "failed to get account id")
		return
	}

	client := &Client{
		hub:    e.hub,
		ws:     ws,
		userId: userId,
		send:   make(chan []byte, 512),
	}
	client.hub.register <- client
	ws.SetCloseHandler(func(code int, text string) error {
		client.hub.unregister <- client
		return nil
	})

	err = e.handler.HandleInit(r.Context(), client)
	if err != nil {
		log.Error().Err(err).
			Str("userId", userId).
			Msg("Closing websocket: failed to initialize websocket")

		closeWs(ws, 5001, "failed to initialize")
		client.hub.unregister <- client

		return
	}

	go client.writer()
	client.reader(r.Context(), e.handler.HandleMessage)
}
