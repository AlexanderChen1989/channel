package client

import "golang.org/x/net/websocket"

// Socket interface for easy test
type Socket interface {
	Send(*Msg) error
	Recv() (*Msg, error)
	Close() error
}

// WSocket Socket implementation for web socket
type WSocket struct {
	conn *websocket.Conn
}

// Send implments Socket.Send
func (ws *WSocket) Send(msg *Msg) error {
	return websocket.JSON.Send(ws.conn, msg)
}

// Close implements Socket.Close
func (ws *WSocket) Close() error {
	return ws.conn.Close()
}

// Recv implements Socket.Recv
func (ws *WSocket) Recv() (*Msg, error) {
	msg := &Msg{}
	return msg, websocket.JSON.Receive(ws.conn, msg)
}
