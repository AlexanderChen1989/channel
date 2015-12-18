package server

import (
	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

type WSHandler interface {
	WSHandle(ctx context.Context, ws *websocket.Conn)
}

type WSHandleFunc func(ctx context.Context, ws *websocket.Conn)

func (fn WSHandleFunc) WSHandle(ctx context.Context, ws *websocket.Conn) {
	fn(ctx, ws)
}

type JoinHandler struct {
	next WSHandler
}

func (h *JoinHandler) WSHandle(ctx context.Context, ws *websocket.Conn) {
}
