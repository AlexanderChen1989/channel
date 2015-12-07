package client

// Message message
type Message struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Ref     string      `json:"ref"`
	Payload interface{} `json:"payload"`
}

// Puller channel to receive result
type Puller struct {
	center *regCenter
	ch     chan *Message
	key    string
}

// Close return puller to regCenter
func (puller *Puller) Close() {
	puller.center.unregister(puller)
}

// Pull pull messge from chan
func (puller *Puller) Pull() <-chan *Message {
	return puller.ch
}
