package server

type client struct {
	topic string
	tr    *Transport
}

type Channel struct {
	clients []client
	onJoin  func(*Message) (map[string]interface{}, error)
}

func (ch *Channel) Push(cli *client, payload interface{}) error {
	return nil
}

func (ch *Channel) Broadcast(cli *client, payload interface{}) error {
	return nil
}

func (ch *Channel) BroadcastToOthers(cli *client, payload interface{}) error {
	return nil
}

func (ch *Channel) Join(tr *Transport, topic string) error {

	return nil
}
