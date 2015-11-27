package channel

type Connection struct{}

func ConnectTo() *Connection {
	conn := &Connection{}
	return conn
}

func (conn *Connection) JoinTo() *Channel {
	ch := &Channel{}
	return ch
}

type Message struct {
	Topic   string                 `json:"topic"`
	Event   string                 `json:"event"`
	Payload map[string]interface{} `json:"payload"`
	Ref     int                    `json:"ref"`
}
