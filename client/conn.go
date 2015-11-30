package client

import (
	"errors"
	"fmt"
	"net/url"
	"path"
	"sync"

	"golang.org/x/net/websocket"
)

const VSN = "1.0.0"

const (
	ConnConnecting = "connecting"
	ConnOpen       = "open"
	ConnClosing    = "closing"
	Connclosed     = "closed"
)

type Conn struct {
	sock    Socket
	status  string
	alive   bool
	mcpool  *sync.Pool
	mpool   *sync.Pool
	bufNum  int
	chM     map[interface{}]map[*MsgCh]bool
	chMLock sync.RWMutex
	counter int
}

func (conn *Conn) loop() {
	conn.alive = true
	for conn.alive {
		msg, err := conn.sock.Recv()
		if err != nil {
			fmt.Printf("%s\n", err)
			break
		}
		conn.dispatch(msg)
	}
	conn.alive = false
}

func (conn *Conn) Recv() *MsgCh {
	return conn.register(&allMsg)
}

func (conn *Conn) Close() error {
	conn.alive = false
	return conn.sock.Close()
}

func (conn *Conn) Send(msg *Msg) error {
	return conn.sock.Send(msg)
}

func (conn *Conn) makeRef() string {
	conn.counter++
	return fmt.Sprint(conn.counter)
}

func ConnectTo(_url string, args url.Values) (*Conn, error) {
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
	bufNum := 2
	conn := &Conn{
		sock:   &WSocket{conn: wconn},
		status: ConnOpen,
		bufNum: bufNum,
		mcpool: &sync.Pool{
			New: func() interface{} {
				return &MsgCh{
					ch: make(chan *Msg, bufNum),
				}
			},
		},
		mpool: &sync.Pool{
			New: func() interface{} {
				return make(map[*MsgCh]bool)
			},
		},
		chM: make(map[interface{}]map[*MsgCh]bool),
	}

	go conn.loop()

	return conn, nil
}

func sendToCh(mch *MsgCh, msg *Msg) {
	select {
	case mch.ch <- msg:
	default:
	}
}

func (conn *Conn) sendToChs(mchs map[*MsgCh]bool, msg *Msg) {
	for mch := range mchs {
		go sendToCh(mch, msg)
	}
}

func (conn *Conn) dispatch(msg *Msg) {
	conn.chMLock.RLock()
	defer conn.chMLock.RUnlock()

	go conn.sendToChs(conn.chM[&allMsg], msg)
	go conn.sendToChs(conn.chM[msg.Topic], msg)
	go conn.sendToChs(conn.chM[msg.Topic+msg.Event], msg)
	go conn.sendToChs(conn.chM[msg.Ref], msg)
}

var allMsg int

func (conn *Conn) register(key interface{}) *MsgCh {
	conn.chMLock.Lock()
	defer conn.chMLock.Unlock()

	ch := conn.mcpool.Get().(*MsgCh)
	ch.conn = conn
	ch.key = key

	// FIXME: ?
	// check whether chan is open and clean elements in open chan
	open := true
	for i := 0; i < conn.bufNum; i++ {
		select {
		case _, open = <-ch.ch:
			if !open {
				break
			}
		default:
			break
		}
	}
	if !open {
		ch.ch = make(chan *Msg, conn.bufNum)
	}

	if conn.chM[key] == nil {
		conn.chM[key] = conn.mpool.Get().(map[*MsgCh]bool)
	}
	conn.chM[key][ch] = true
	return ch
}

func (conn *Conn) unregister(key interface{}, ch *MsgCh) {
	conn.chMLock.Lock()
	defer conn.chMLock.Unlock()

	m := conn.chM[key]
	if m != nil {
		delete(m, ch)
		conn.mcpool.Put(ch)
	}
	if len(m) == 0 {
		conn.mpool.Put(m)
	}
}
