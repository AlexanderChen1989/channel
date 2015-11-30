package client

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChan(t *testing.T) {
	timeout := 10 * time.Millisecond
	conn := NewConn(2, timeout)
	conn.status = ConnOpen
	sock := &socket{timeout: timeout}
	sock.msg = &Msg{
		Topic:   "rooms:alex",
		Event:   "new:msg",
		Payload: "",
	}
	conn.sock = sock
	ch1, err := conn.Chan("rooms:alex")
	assert.Nil(t, err)
	ch2, err := conn.Chan("rooms:aaron")
	assert.Nil(t, err)
	for i := 0; i < 20; i++ {
		select {
		case <-ch2.Recv().Recv():
			t.Error("Should not recv msg.")
		default:
		}
		select {
		case <-ch1.Recv().Recv():
		case <-time.After(timeout + 100*time.Millisecond):
			t.Error("Should recv msg.")
		}
	}

	ch, err := conn.Chan("rooms:alex")

	for i := 0; i < 10; i++ {
		assert.Nil(t, err)
		mch, err := ch.Request("new:msg", "")
		assert.Nil(t, err)
		select {
		case msg := <-mch.Recv():
			fmt.Println(msg)
		case <-time.After(timeout + 100*time.Millisecond):
			t.Error("Shoud recev reply")
		}
		mch.Close()
	}
}

func TestChanReal(t *testing.T) {
	conn, err := ConnectTo("http://localhost:4000/socket", nil)
	assert.Nil(t, err)
	defer conn.Close()
	ch, err := conn.Chan("rooms:lobby")
	assert.Nil(t, err)
	mch, err := ch.Join()
	assert.Nil(t, err)
	defer mch.Close()
}
