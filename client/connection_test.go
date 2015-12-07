package client

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"golang.org/x/net/context"
)

func pullMessage(puller *Puller, timeout time.Duration, errCh chan error) {
	select {
	case <-puller.Pull():
		errCh <- nil
	case <-time.After(timeout):
		errCh <- errors.New("Timeout")
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

	const num = 100
	for i := 0; i < num; i++ {
		msg := &Message{
			Topic:   "topic" + fmt.Sprint(num%3),
			Event:   "new_msg" + fmt.Sprint(num%10),
			Ref:     conn.ref.makeRef(),
			Payload: "msg",
		}
		all := conn.OnMessage()
		errCh := make(chan error)
		go pullMessage(all, 5*time.Millisecond, errCh)
		conn.push(msg)
		assert.Nil(t, <-errCh)

		ch1, err := conn.Chan(msg.Topic)
		assert.Nil(t, err)
		ch2, err := conn.Chan(msg.Topic + "xx")
		assert.Nil(t, err)

		pullerCh1 := ch1.OnMessage()
		pullerCh2 := ch2.OnMessage()
		pullerCh1ErrCh := make(chan error)
		pullerCh2ErrCh := make(chan error)
		go pullMessage(pullerCh1, 5*time.Millisecond, pullerCh1ErrCh)
		go pullMessage(pullerCh2, 5*time.Millisecond, pullerCh2ErrCh)
		conn.push(msg)
		assert.Nil(t, <-pullerCh1ErrCh)
		assert.NotNil(t, <-pullerCh2ErrCh)

		pullerEvt1 := ch1.OnEvent(msg.Event)
		pullerEvt2 := ch1.OnEvent(msg.Event + "xx")
		pullerEvt1ErrCh := make(chan error)
		pullerEvt2ErrCh := make(chan error)
		go pullMessage(pullerEvt1, 5*time.Millisecond, pullerEvt1ErrCh)
		go pullMessage(pullerEvt2, 5*time.Millisecond, pullerEvt2ErrCh)
		conn.push(msg)
		assert.Nil(t, <-pullerEvt1ErrCh)
		assert.NotNil(t, <-pullerEvt2ErrCh)

		req, err := ch1.Request("hello_world", "body")
		assert.Nil(t, err)
		reqErrCh := make(chan error)
		go pullMessage(req, 5*time.Millisecond, reqErrCh)
	}
}
