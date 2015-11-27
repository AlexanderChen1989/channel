package channel

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"sync"

	"golang.org/x/net/context"
	"golang.org/x/net/websocket"
)

var cmd = `
const VSN = "1.0.0"
const SOCKET_STATES = {connecting: 0, open: 1, closing: 2, closed: 3}
const CHANNEL_STATES = {
  closed: "closed",
  errored: "errored",
  joined: "joined",
  joining: "joining",
}
const CHANNEL_EVENTS = {
  close: "phx_close",
  error: "phx_error",
  join: "phx_join",
  reply: "phx_reply",
  leave: "phx_leave"
}
const TRANSPORTS = {
  longpoll: "longpoll",
  websocket: "websocket"
}
`

const (
	VSN = "1.0.0"

	ChannelClosed  = "closed"
	ChannelErrored = "errored"
	ChannelJoined  = "joined"
	ChannelJoining = "joining"

	EventPhxClose = "phx_close"
	EventPhxError = "phx_error"
	EventPhxJoin  = "phx_join"
	EventPhxReply = "phx_reply"
	EventPhxLeave = "phx_leave"
)

type Connection struct {
	originURL string
	socketURL string
	lock      sync.Mutex
	ctx       context.Context
	conn      *websocket.Conn
	chans     map[string]*Channel
	count     int
}

func (conn *Connection) connect() (err error) {
	conn.conn, err = websocket.Dial(conn.socketURL, "", conn.originURL)
	return
}

func (conn *Connection) Close() error {
	return conn.conn.Close()
}

func (conn *Connection) Ref() int {
	// FIXME: thread safe!
	conn.count++
	return conn.count
}

func ConnectTo(_url string, args url.Values) (*Connection, error) {
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

	conn := &Connection{
		ctx:       context.Background(),
		chans:     make(map[string]*Channel),
		originURL: fmt.Sprintf("%s://%s", oscheme, surl.Host),
		socketURL: surl.String(),
	}

	return conn, conn.connect()
}

func (conn *Connection) JoinTo(topic string) (*Channel, error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	topic = strings.TrimSpace(topic)
	if topic == "" {
		return nil, errors.New("Topic is empty.")
	}

	if conn.chans[topic] != nil {
		return nil, errors.New("Allready joined to '" + topic + "'.")
	}

	ctx, cancel := context.WithCancel(conn.ctx)

	ch := &Channel{
		ctx:    ctx,
		conn:   conn,
		topic:  topic,
		cancel: cancel,
		msgCh:  make(chan *Message),
	}

	conn.chans[topic] = ch

	return ch, nil
}
