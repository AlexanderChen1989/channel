package client

// Puller channel to receive result
type Puller struct {
	center *regCenter
	ch     chan *Message
	key    string
}

// Close return puller to regCenter
func (puller *Puller) Close() {
	if puller.center == nil {
		return
	}
	puller.center.unregister(puller)
}

// Pull pull messge from chan
func (puller *Puller) Pull() <-chan *Message {
	return puller.ch
}
