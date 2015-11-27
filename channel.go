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
	ctx    context.Context
	conn   *Connection
	topic  string
	status string
	cancel func()
	msgCh  chan *Message

	refMaplock    sync.RWMutex
	refMap        map[int]chan *Message
	recvChslock   sync.RWMutex
	recvChs       map[string]chan *Message
	recvAllChlock sync.Mutex
	recvAllCh     chan *Message
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

func (ch *Channel) dispatchToRefMap(msg *Message) {
	rch := ch.refMagGet(msg.Ref)
	if rch == nil {
		return
	}
	select {
	case rch <- msg:
	case <-ch.ctx.Done():
		return
	default:
	}
	delete(ch.refMap, msg.Ref)
}

func (ch *Channel) dispatchToRecvChs(msg *Message) {
	rch := ch.recvChsGet(msg.Event)
	if rch == nil {
		return
	}
	select {
	case rch <- msg:
	case <-ch.ctx.Done():
		return
	default:
	}
}

func (ch *Channel) dispatchToRecvAllCh(msg *Message) {
	if ch.recvAllCh == nil {
		return
	}
	select {
	case ch.recvAllCh <- msg:
	case <-ch.ctx.Done():
		return
	default:
	}
}

func (ch *Channel) dispatch(msg *Message) {
	ch.dispatchToRefMap(msg)
	ch.dispatchToRecvChs(msg)
	ch.dispatchToRecvAllCh(msg)
}

func newChannel(conn *Connection, topic string) *Channel {
	ctx, cancel := context.WithCancel(conn.ctx)

	ch := &Channel{
		ctx:     ctx,
		conn:    conn,
		topic:   topic,
		status:  ChannelReady,
		cancel:  cancel,
		msgCh:   make(chan *Message),
		refMap:  make(map[int]chan *Message),
		recvChs: make(map[string]chan *Message),
	}

	return ch
}

func (ch *Channel) refMagGet(ref int) chan *Message {
	ch.refMaplock.RLock()
	defer ch.refMaplock.RUnlock()

	return ch.refMap[ref]
}

func (ch *Channel) refMapAdd(ref int, mch chan *Message) {
	ch.refMaplock.Lock()
	defer ch.refMaplock.Unlock()

	ch.refMap[ref] = mch
}

func (ch *Channel) Request(msg *Message) (chan *Message, error) {
	msg.Ref = ch.conn.Ref()
	msg.Topic = ch.topic

	if err := ch.conn.Send(msg); err != nil {
		return nil, err
	}

	mch := make(chan *Message, 1)

	ch.refMapAdd(msg.Ref, mch)

	return mch, nil
}

var ErrRecvChTaken = errors.New("Recv chan has been taken.")

func (ch *Channel) recvChsGet(event string) chan *Message {
	ch.recvChslock.RLock()
	defer ch.recvChslock.RUnlock()

	return ch.recvChs[event]
}

func (ch *Channel) recvChsAdd(event string, mch chan *Message) {
	ch.recvChslock.Lock()
	defer ch.recvChslock.Unlock()

	ch.recvChs[event] = mch
}

func (ch *Channel) Recv(event string) (chan *Message, error) {
	if ch.recvChsGet(event) != nil {
		return nil, ErrRecvChTaken
	}
	mch := make(chan *Message, 2)
	ch.recvChsAdd(event, mch)
	return mch, nil
}

func (ch *Channel) RecvAll() (chan *Message, error) {
	ch.recvAllChlock.Lock()
	defer ch.recvAllChlock.Unlock()

	if ch.recvAllCh != nil {
		return nil, ErrRecvChTaken
	}

	ch.recvAllCh = make(chan *Message, 2)
	return ch.recvAllCh, nil
}
