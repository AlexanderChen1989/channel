package client

type Chan struct {
	Parent
	Topic string
	msgCh chan *Msg
	pullM map[interface{}]chan *Msg
}

func NewChan(parent Parent, topic string) *Chan {
	ch := &Chan{Parent: parent}
	ch.Topic = topic
	ch.pullM = map[interface{}]chan *Msg{}
	return ch
}

// PushPullRemover
func (ch *Chan) Pull(key interface{}, mch chan *Msg) error {
	ch.pullM[key] = mch
	return nil
}

func (ch *Chan) Push(msg *Msg) error {
	msg.Topic = ch.Topic
	ch.Parent.Push(msg)
	return nil
}

func (ch *Chan) Remove(key interface{}) {
	delete(ch.pullM, key)
}

func (ch *Chan) loop() {
	ch.Parent.Pull(ch.Topic, ch.msgCh)
	for {
		select {
		case msg := <-ch.msgCh:
			go ch.dispatch(msg)
		}
	}
}

func (ch *Chan) dispatch(msg *Msg) error {
	select {
	case ch.pullM[msg.Topic] <- msg:
	default:
	}
	return nil
}

func (ch *Chan) Close() error {
	// FIXME
	ch.Parent.Remove(ch.Topic)
	return nil
}

func (ch *Chan) Request(payload interface{}) *Request {
	msg := &Msg{Payload: payload}
	return NewRequest(ch, msg)
}

func (ch *Chan) Join() error {
	// TODO

	return nil
}
