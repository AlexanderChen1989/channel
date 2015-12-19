package server

type Client struct {
	tr     *Transport
	topic  string
	pullCh chan *Message
	pushCh chan *Message
}

func NewClient(tr *Transport, topic string) *Client {
	return &Client{
		tr:     tr,
		topic:  topic,
		pullCh: make(chan *Message),
		pushCh: tr.pushCh,
	}
}

type Channel struct {
	clients []*Client
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

func (ch *Channel) Join(cli *Client, topic, tail string,
	payload map[string]interface{}) {
}
