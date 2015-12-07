package client

import "time"

// This file is just for test

type socket struct {
	counter int
	timeout time.Duration
	msg     *Message
	msgs    []*Message
}

func (sock *socket) Send(msg *Message) error {
	return nil
}

func (sock *socket) Close() error {
	return nil
}

func (sock *socket) Recv() (*Message, error) {
	return nil, nil
}
