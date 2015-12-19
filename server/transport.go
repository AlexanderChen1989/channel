package server

import (
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type Finder interface {
	Find(pattern string) (*Channel, string)
}

type Transport struct {
	finder Finder
	ctx    context.Context
	conn   *websocket.Conn
	id     string

	clients []*Client
	pushCh  chan *Message
}

func (tr *Transport) pushLoop() {
	for msg := range tr.pushCh {
		websocket.JSON.Send(tr.conn, msg)
	}
}

func (tr *Transport) pullLoop() {
	for {
		msg := &Message{}
		err := websocket.JSON.Receive(tr.conn, &msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		select {
		case <-tr.ctx.Done():
			return
		}

	}
}

func (tr *Transport) Start() {
	// start transport loops
	// start message dispatch loop
}

func NewTransport(finder Finder, ctx context.Context, conn *websocket.Conn, id string) *Transport {
	return &Transport{finder: finder, ctx: ctx, conn: conn, id: id}
}
