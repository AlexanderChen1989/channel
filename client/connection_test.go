package client

import (
	"fmt"
	"testing"

	"golang.org/x/net/context"
)

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

	const num = 1000

	for i := 0; i < num; i++ {
		msg := &Message{
			Topic:   "topic" + fmt.Sprint(num%3),
			Event:   "new_msg" + fmt.Sprint(num%10),
			Ref:     conn.ref.makeRef(),
			Payload: "msg",
		}
		all := conn.OnMessage()
		msg := <-all.Pull()
		all.Close()
	}
}
