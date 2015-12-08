package client

import (
	"errors"
	"testing"

	"golang.org/x/net/context"
)

func pullMessage(puller *Puller) error {
	select {
	case <-puller.Pull():
		return nil
	default:
		return errors.New("no msg")
	}
}

func TestConnection(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	conn := &Connection{
		ctx:    ctx,
		cancel: cancel,
		sock:   newSocket(1),
		center: newRegCenter(),
		msgs:   make(chan *Message),
		status: ConnOpen,
	}
	defer conn.Close()
	conn.start()
	// TODO: add test!
}
