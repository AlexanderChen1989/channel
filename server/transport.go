package server

import (
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type Finder interface {
	Find(pattern string) (*Channel, string)
}

type Transport struct {
	finder  Finder
	ctx     context.Context
	conn    *websocket.Conn
	id      string
	ReadCh  <-chan *Message
	WriteCh chan<- *Message
}

func NewTransport(finder Finder, ctx context.Context, conn *websocket.Conn, id string) *Transport {
	return &Transport{finder: finder, ctx: ctx, conn: conn, id: id}
}
