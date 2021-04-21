package internal

import "encoding/json"

type IncomingMessage struct {
	SenderId string `json:"sender_id"`
	Body     string `json:"body"`
}

func (m *IncomingMessage) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *IncomingMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}


type OutgoingMessage struct {
	ReceiverId string `json:"receiver_id"`
	Body     string `json:"body"`
}

func (m *OutgoingMessage) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

func (m *OutgoingMessage) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
