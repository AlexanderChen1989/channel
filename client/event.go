package client

type EventQueue struct {
	PushPullRemover
	GetPuter
	RefMaker
	Event string
	msgCh chan *Msg
	pullM map[interface{}]chan *Msg
}

// PushPullRemover
func (evtq *EventQueue) Pull(key interface{}, mch chan *Msg) error {
	evtq.pullM[key] = mch
	return nil
}

func (evtq *EventQueue) Push(msg *Msg) error {
	msg.Event = evtq.Event
	evtq.PushPullRemover.Push(msg)
	return nil
}

func (evtq *EventQueue) Remove(key interface{}) {
	delete(evtq.pullM, key)
}

func (evtq *EventQueue) loop() {
	evtq.PushPullRemover.Pull(evtq.Event, evtq.msgCh)
	for {
		select {
		case msg := <-evtq.msgCh:
			go evtq.dispatch(msg)
		}
	}
}

func (evtq *EventQueue) dispatch(msg *Msg) error {
	select {
	case evtq.pullM[msg.Topic] <- msg:
	default:
	}
	return nil
}

func (evtq *EventQueue) Request(payload interface{}) (*Request, error) {
	msg := &Msg{
		Payload: payload,
	}
	// FIXME
	return &Request{Msg: msg}, nil
}
