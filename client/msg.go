package client

type Msg struct {
	Topic   string      `json:"topic"`
	Event   string      `json:"event"`
	Ref     string      `json:"ref"`
	Payload interface{} `json:"payload"`
}

type MsgCh struct {
	conn *Conn
	ch   chan *Msg
	key  interface{}
}

func (mc *MsgCh) Close() {
	if mc.conn != nil {
		mc.conn.unregister(mc.key, mc)
	}
}

func (mc *MsgCh) Recv() <-chan *Msg {
	return mc.ch
}
