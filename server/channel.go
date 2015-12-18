package server

type Client struct {
	tr     *Transport
	topic  string
	pullCh chan *Message
	pushCh chan *Message
}

type Channel struct {
	clients []Client
	onJoin  func(payload map[string]interface{}) (map[string]interface{}, error)
}

func (ch *Channel) Push(cli *Client, payload interface{}) error {
	return nil
}

func (ch *Channel) Broadcast(cli *Client, payload interface{}) error {
	return nil
}

func (ch *Channel) BroadcastToOthers(cli *Client, payload interface{}) error {
	return nil
}

func (ch *Channel) Join(tr *Transport, topic, tail string,
	payload map[string]interface{}) {
}
