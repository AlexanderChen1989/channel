package client

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	ConnConnecting = "connecting"
	ConnOpen       = "open"
	ConnClosing    = "closing"
	ConnClosed     = "closed"
)

func Connect(_url string, args url.Values) (*Conn, error) {
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

	conn := &Conn{
		status:    ConnConnecting,
		originURL: fmt.Sprintf("%s://%s", oscheme, surl.Host),
		socketURL: surl.String(),
	}
	if err := conn.connect(); err != nil {
		return nil, err
	}
	conn.status = ConnOpen
	go conn.loop()

	return conn, nil
}

type Conn struct {
	*websocket.Conn
	status    string
	originURL string
	socketURL string
	counter   int
	pullM     map[interface{}]chan *Msg
}

func (conn *Conn) loop() {
	for {
		var msg Msg
		// FIXME: ?
		err := websocket.JSON.Receive(conn.Conn, &msg)
		if err != nil {
			fmt.Println(err)
		}
		conn.dispatch(&msg)
	}
}

func (conn *Conn) connect() (err error) {
	conn.Conn, err = websocket.Dial(conn.socketURL, "", conn.originURL)
	return
}

func (conn *Conn) dispatch(msg *Msg) error {
	select {
	case conn.pullM[msg.Topic] <- msg:
	default:
	}
	return nil
}

// PushPullRemover
func (conn *Conn) Push(msg *Msg) error {
	return websocket.JSON.Send(conn.Conn, msg)
}

func (conn *Conn) Pull(key interface{}, ch chan *Msg) error {
	conn.pullM[key] = ch
	return nil
}

func (conn *Conn) Remove(key interface{}) {
	delete(conn.pullM, key)
}

// Closer
func (conn *Conn) Close() error {
	return conn.Conn.Close()
}

// GetPuter
func (conn *Conn) Get() chan *Msg {
	return make(chan *Msg, 1)
}
func (conn *Conn) Put(chan *Msg) {}

// RefMaker
func (conn *Conn) MakeRef() string {
	// FIXME
	conn.counter++
	return fmt.Sprint(conn.counter)
}

func (conn *Conn) Chan(topic string) (*Chan, error) {
	topic = strings.TrimSpace(topic)
	ch := NewChan(conn, topic)
	if err := ch.Join(); err != nil {
		ch.Close()
		return nil, err
	}
	return ch, nil
}
