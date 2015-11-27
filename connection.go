package channel

import (
	"errors"
	"fmt"
	"log"
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
	cancel    func()
	status    string
	conn      *websocket.Conn
	sendLock  sync.Mutex
	chans     map[string]*Channel
	count     int
}

func (conn *Connection) connect() (err error) {
	conn.conn, err = websocket.Dial(conn.socketURL, "", conn.originURL)
	return
}

func (conn *Connection) Close() error {
	conn.cancel()
	return conn.conn.Close()
}

func (conn *Connection) Ref() int {
	// FIXME: thread safe!
	conn.count++
	return conn.count
}

func (conn *Connection) Send(msg *Message) error {
	conn.sendLock.Lock()
	defer conn.sendLock.Unlock()

	return websocket.JSON.Send(conn.conn, msg)
}

func (conn *Connection) loop() {
	// FIXME: ?
	for {
		var msg Message
		err := websocket.JSON.Receive(conn.conn, &msg)
		if err != nil {
			log.Println(err)
			continue
		}
		log.Println(msg)
		if conn.chans[msg.Topic] == nil {
			continue
		}
		select {
		case conn.chans[msg.Topic].msgCh <- &msg:
		default:
		}
	}

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

	ctx, cancel := context.WithCancel(context.Background())
	conn := &Connection{
		ctx:       ctx,
		cancel:    cancel,
		chans:     make(map[string]*Channel),
		originURL: fmt.Sprintf("%s://%s", oscheme, surl.Host),
		socketURL: surl.String(),
	}
	if err := conn.connect(); err != nil {
		return nil, err
	}
	go conn.loop()

	return conn, nil
}

func (conn *Connection) removeChan(ch *Channel) {
	ch.Close()
	delete(conn.chans, ch.topic)
}

func (conn *Connection) addChan(ch *Channel) {
	conn.chans[ch.topic] = ch
	go ch.loop()
}

func (conn *Connection) checkTopic(topic string) error {
	if topic == "" {
		return errors.New("Topic is empty.")
	}

	if conn.chans[topic] != nil {
		return errors.New("Allready joined to '" + topic + "'.")
	}

	return nil
}

func (conn *Connection) JoinTo(topic string) (*Channel, error) {
	conn.lock.Lock()
	defer conn.lock.Unlock()

	topic = strings.TrimSpace(topic)
	if err := conn.checkTopic(topic); err != nil {
		return nil, err
	}

	ch := NewChannel(conn, topic)
	conn.addChan(ch)

	if err := ch.Join(); err != nil {
		conn.removeChan(ch)
		return nil, err
	}

	return ch, nil
}
