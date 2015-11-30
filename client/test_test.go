package client

import (
	"fmt"
	"sync"
	"time"
)

// This file is just for test

type socket struct {
	counter int
	timeout time.Duration
	msg     *Msg
	msgs    []*Msg
}

func (sock *socket) Send(msg *Msg) error {
	msg.Payload = "reply for " + msg.Ref
	sock.msgs = append(sock.msgs, msg)
	return nil
}

func (sock *socket) Close() error {
	return nil
}

func (sock *socket) Recv() (*Msg, error) {
	time.Sleep(sock.timeout)
	if len(sock.msgs) > 0 {
		msg := sock.msgs[len(sock.msgs)-1]
		sock.msgs = sock.msgs[:len(sock.msgs)-1]
		return msg, nil
	}
	sock.counter++
	if sock.msg != nil {
		sock.msg.Ref = fmt.Sprint(sock.counter)
		return sock.msg, nil
	}
	return &Msg{
		Topic:   "rooms:alex",
		Event:   "new:msg",
		Ref:     fmt.Sprint(sock.counter),
		Payload: "Hello",
	}, nil
}

func NewConn(bufNum int, timeout time.Duration) *Conn {
	conn := &Conn{
		sock:   &socket{timeout: timeout},
		status: ConnOpen,
		bufNum: bufNum,
		mcpool: &sync.Pool{
			New: func() interface{} {
				return &MsgCh{
					ch: make(chan *Msg, bufNum),
				}
			},
		},
		mpool: &sync.Pool{
			New: func() interface{} {
				return make(map[*MsgCh]bool)
			},
		},
		chM: make(map[interface{}]map[*MsgCh]bool),
	}

	go conn.loop()

	return conn
}
