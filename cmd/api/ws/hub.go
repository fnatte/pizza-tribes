package ws

type Message struct {
	Recipient string
	Body      []byte
}

type Hub struct {
	clients    map[*Client]bool
	message    chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		message:    make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) userIds() map[string]bool {
	var accIds = map[string]bool{}
	for client := range h.clients {
		accIds[client.userId] = true
	}
	return accIds
}

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
		case msg := <-h.message:
			for client := range h.clients {
				if client.userId == msg.Recipient {
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
