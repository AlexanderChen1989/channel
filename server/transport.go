package server

import (
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type Transport struct {
	ctx     context.Context
	conn    *websocket.Conn
	ReadCh  <-chan *Message
	WriteCh chan<- *Message
}
