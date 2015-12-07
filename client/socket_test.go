package client

import "time"

// This file is just for test

type socket struct {
	timeout time.Duration
	ref     refMaker
	msgs    chan *Message
}

func (sock *socket) Send(msg *Message) error {
	sock.msgs <- msg
	return nil
}

func (sock *socket) Close() error {
	return nil
}

func (sock *socket) Recv() (*Message, error) {
	time.Sleep(100 * time.Millisecond)
	msg := <-sock.msgs
	switch payload := msg.Payload.(type) {
	case string:
		msg.Payload = "reply for " + payload
	}
	return msg, nil
}

func newSocket(size int) *socket {
	return &socket{
		msgs: make(chan *Message, size),
	}
}
