package client

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	const num = 100
	for i := 0; i < num; i++ {
		msg := &Message{
			Topic:   "topic" + fmt.Sprint(num%3),
			Event:   "new_msg" + fmt.Sprint(num%10),
			Ref:     conn.ref.makeRef(),
			Payload: "msg",
		}
		all := conn.OnMessage()
		ch1, _ := conn.Chan(msg.Topic)
		ch2, _ := conn.Chan(msg.Topic + "xx")
		pullerCh1 := ch1.OnMessage()
		pullerCh2 := ch2.OnMessage()
		pullerEvt1 := ch1.OnEvent(msg.Event)
		pullerEvt2 := ch1.OnEvent(msg.Event + "xx")
		req, _ := ch1.Request("hello_world", "body")

		conn.msgs <- msg
		time.Sleep(5 * time.Millisecond)

		assert.Nil(t, pullMessage(all))
		assert.Nil(t, pullMessage(pullerCh1))
		assert.NotNil(t, pullMessage(pullerCh2))
		assert.Nil(t, pullMessage(pullerEvt1))
		assert.NotNil(t, pullMessage(pullerEvt2))
		assert.Nil(t, pullMessage(req))
	}
}
