package ws

type Message struct {
	Recipient string
	Body      []byte
}

// The Hub maintain all web socket clients. When started (using Run()),
// the Hub will pass messages from the message channel to the corresponding
// recipient. The SendTo() func can be used to place messages on the message
// channel.
type Hub struct {
	clients    map[*Client]bool
	messages   chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		messages:   make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// Sends bytes to a recipient (user id) by putting a message on the messages
// channel. If the Hub is running (using Run()), the messages will be passed
// to the corresponding Client.
func (h *Hub) SendTo(recipient string, body []byte) {
	h.messages <- &Message{Recipient: recipient, Body: body}
}

// Starts the Hub (and blocks). It will pull messages and send them to the
// corresponding client.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case msg := <-h.messages:
			// TODO: should probably optimize this
			for client := range h.clients {
				if client.userId == msg.Recipient || msg.Recipient == "everyone" {
					select {
					case client.send <- []byte(msg.Body):
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			}
		}
	}
}
