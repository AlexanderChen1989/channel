package client

// Message message
type Message struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Ref     string      `json:"ref"`
	Payload interface{} `json:"payload"`
}
