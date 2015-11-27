package channel

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
)

type Channel struct {
	ctx  context.Context
	conn *Connection

	topic  string
	cancel func()
	msgCh  chan *Message

	counter   int
	lock      sync.Mutex
	recvChs   map[string]chan *Message
	recvAllCh chan *Message
	refMap    map[int]chan *Message
}

func (ch *Channel) loop() {
	for {
		select {
		case <-ch.ctx.Done():
			break
		case msg := <-ch.msgCh:
			go ch.dispatch(msg)
		}
	}
}

func (ch *Channel) dispatch(msg *Message) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	select {
	case ch.refMap[msg.Ref] <- msg:
	default:
	}
	delete(ch.refMap, msg.Ref)

	select {
	case ch.recvChs[msg.Event] <- msg:
	default:
	}

	select {
	case ch.recvAllCh <- msg:
	default:
	}
}

func Msg(event string, payload interface{}) *Message {
	return nil
}

func (ch *Channel) mkRef() int {
	return 0
}

func (ch *Channel) Request(msg *Message) (chan *Message, error) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if ch.refMap == nil {
		ch.refMap = make(map[int]chan *Message)
	}

	msg.Ref = ch.conn.Ref()
	msg.Topic = ch.topic
	mch := make(chan *Message, 1)

	ch.refMap[msg.Ref] = mch

	return mch, nil
}

func (ch *Channel) Send(msg *Message) error {
	return nil
}

var ErrRecvChTaken = errors.New("Recv chan has been taken.")

func (ch *Channel) Recv(event string) (chan *Message, error) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if ch.recvChs == nil {
		ch.recvChs = make(map[string]chan *Message)
	}

	if ch.recvChs[event] != nil {
		return nil, ErrRecvChTaken
	}

	mch := make(chan *Message, 2)
	ch.recvChs[event] = mch
	return mch, nil
}

func (ch *Channel) RecvAll() (chan *Message, error) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	if ch.recvAllCh != nil {
		return nil, ErrRecvChTaken
	}

	ch.recvAllCh = make(chan *Message, 2)
	return ch.recvAllCh, nil
}
