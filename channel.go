package channel

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"golang.org/x/net/context"
)

const (
	ChannelClosed  = "closed"
	ChannelErrored = "errored"
	ChannelJoined  = "joined"
	ChannelJoining = "joining"
	ChannelReady   = "ready"
)

type Channel struct {
	ctx  context.Context
	conn *Connection

	topic  string
	status string
	cancel func()
	msgCh  chan *Message

	lock      sync.Mutex
	recvChs   map[string]chan *Message
	recvAllCh chan *Message
	refMap    map[int]chan *Message
}

func (ch *Channel) Close() {
	ch.cancel()
}

func (ch *Channel) loop() {
	for {
		select {
		case <-ch.ctx.Done():
			return
		case msg := <-ch.msgCh:
			go ch.dispatch(msg)
		}
	}
}

func (ch *Channel) Join() error {
	if ch.status != ChannelReady {
		return errors.New("Channel for " + ch.topic + " status " + ch.status + ".")
	}
	msg := &Message{
		Topic:   ch.topic,
		Event:   EventPhxJoin,
		Payload: "",
		Ref:     ch.conn.Ref(),
	}
	mch, err := ch.Request(msg)
	if err != nil {
		return err
	}

	mchall, err := ch.RecvAll()

	select {
	case reply := <-mch:
		fmt.Println(reply)
		return nil
	case reply := <-mchall:
		fmt.Println(reply)
		return nil
	case <-time.After(5 * time.Second):
		return errors.New("Join timeout.")
	}
}

func (ch *Channel) dispatch(msg *Message) {
	ch.lock.Lock()
	defer ch.lock.Unlock()

	select {
	case ch.refMap[msg.Ref] <- msg:
	case <-ch.ctx.Done():
		return
	default:
	}
	delete(ch.refMap, msg.Ref)

	select {
	case ch.recvChs[msg.Event] <- msg:
	case <-ch.ctx.Done():
		return
	default:
	}

	select {
	case ch.recvAllCh <- msg:
	case <-ch.ctx.Done():
		return
	default:
	}
}

func (ch *Channel) Request(msg *Message) (chan *Message, error) {
	msg.Ref = ch.conn.Ref()
	msg.Topic = ch.topic
	if err := ch.conn.Send(msg); err != nil {
		return nil, err
	}

	mch := make(chan *Message, 1)

	ch.lock.Lock()
	defer ch.lock.Unlock()

	ch.refMap[msg.Ref] = mch

	return mch, nil
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
