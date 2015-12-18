package server

type Channel struct {
	trans  []Transport
	onJoin func(*Message) (map[string]interface{}, error)
}

func (ch *Channel) Push(tr *Transport, payload interface{}) error {
	return nil
}

func (ch *Channel) Broadcast(tr *Transport, payload interface{}) error {
	return nil
}

func (ch *Channel) BroadcastToOthers(tr *Transport, payload interface{}) error {
	return nil
}

func (ch *Channel) Join(tr *Transport) error {

	return nil
}
