package client

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"sync"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

const VSN = "1.0.0"

const (
	ConnConnecting = "connecting"
	ConnOpen       = "open"
	ConnClosing    = "closing"
	Connclosed     = "closed"
)

type Connection struct {
	ctx    context.Context
	cancel func()

	sock   Socket
	ref    refMaker
	center *regCenter
	msgs   chan *Message

	status string
}

func Connect(_url string, args url.Values) (*Connection, error) {
	surl, err := url.Parse(_url)

	if err != nil {
		return nil, err
	}

	if !surl.IsAbs() {
		return nil, errors.New("URL should be absolute.")
	}

	oscheme := surl.Scheme
	switch oscheme {
	case "http":
		surl.Scheme = "ws"
	case "https":
		surl.Scheme = "wss"
	default:
		return nil, errors.New("Schema should be http or https.")
	}

	surl.Path = path.Join(surl.Path, "websocket")
	surl.RawQuery = args.Encode()

	originURL := fmt.Sprintf("%s://%s", oscheme, surl.Host)
	socketURL := surl.String()

	wconn, err := websocket.Dial(socketURL, "", originURL)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	conn := &Connection{
		ctx:    ctx,
		cancel: cancel,
		sock:   &WSocket{conn: wconn},
		center: newRegCenter(),
		msgs:   make(chan *Message),
		status: ConnOpen,
	}

	conn.start()

	return conn, nil
}

const all = ""

// OnMessage receive all message on connection
func (conn *Connection) OnMessage() *Puller {
	return conn.center.register(all)
}

func (conn *Connection) push(msg *Message) error {
	return conn.sock.Send(msg)
}

func (conn *Connection) heartbeatLoop() {
	msg := &Message{
		Topic:   "phoenix",
		Event:   "heartbeat",
		Ref:     conn.ref.makeRef(),
		Payload: "",
	}
	for {
		select {
		case <-time.After(30000 * time.Millisecond):
			conn.sock.Send(msg)
		case <-conn.ctx.Done():
			return

		}

	}
}

func (conn *Connection) pullLoop() {
	for {
		msg, err := conn.sock.Recv()
		if err != nil {
			fmt.Printf("%s\n", err)
			close(conn.msgs)
			return
		}
		select {
		case <-conn.ctx.Done():
			return
		case conn.msgs <- msg:
		}
	}
}

func (conn *Connection) coreLoop() {
	for {
		select {
		case <-conn.ctx.Done():
			return
		case msg, ok := <-conn.msgs:
			if !ok {
				return
			}
			conn.dispatch(msg)
		}
	}
}

func (conn *Connection) start() {
	go conn.pullLoop()
	go conn.coreLoop()
	go conn.heartbeatLoop()
}

func (conn *Connection) Close() error {
	conn.cancel()
	return conn.sock.Close()
}

func (conn *Connection) pushToChan(puller *Puller, msg *Message) {
	select {
	case puller.ch <- msg:
	case <-time.After(10 * time.Millisecond):
	}
}

func (conn *Connection) pushToChans(wg *sync.WaitGroup, pullers []*Puller, msg *Message) {
	for _, puller := range pullers {
		go conn.pushToChan(puller, msg)
	}
	wg.Done()
}

func (conn *Connection) dispatch(msg *Message) {
	var wg sync.WaitGroup
	wg.Add(4)
	go conn.pushToChans(&wg, conn.center.getPullers(all), msg)
	go conn.pushToChans(&wg, conn.center.getPullers(toKey(msg.Topic, "", "")), msg)
	go conn.pushToChans(&wg, conn.center.getPullers(toKey(msg.Topic, msg.Event, "")), msg)
	go conn.pushToChans(&wg, conn.center.getPullers(toKey("", "", msg.Ref)), msg)
	wg.Wait()
}
