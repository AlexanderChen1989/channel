package channel

type Channel struct {
	msgCh   chan *Message
	filters map[string][]chan *Message
}

func (ch *Channel) RequestTo(event string, payload map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (ch *Channel) Send(event string, payload map[string]interface{}) (*Message, error) {
	return nil, nil
}

func (ch *Channel) Recv(event string) (*Message, error) {
	return nil, nil
}

func (ch *Channel) RecvAll() (*Message, error) {
	return nil, nil
}
