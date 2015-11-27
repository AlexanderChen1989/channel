package channel

import (
	"errors"
	"fmt"
	"net/url"
	"path"

	"golang.org/x/net/websocket"
)

type Connection struct {
	originURL string
	socketURL string
	conn      *websocket.Conn
}

func (conn *Connection) connect() (err error) {
	conn.conn, err = websocket.Dial(conn.socketURL, "", conn.originURL)
	return
}

func (conn *Connection) Close() error {
	return conn.conn.Close()
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
		originURL: fmt.Sprintf("%s://%s", oscheme, surl.Host),
		socketURL: surl.String(),
	}

	return conn, conn.connect()
}

func (conn *Connection) JoinTo(topic string) *Channel {
	ch := &Channel{}
	return ch
}

func (conn *Connection) getCh() chan map[string]interface{} {
	return nil
}

func (conn *Connection) putCh(ch chan map[string]interface{}) {

}
