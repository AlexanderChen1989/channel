package client

import "errors"

const (
	ChClosed  = "closed"
	ChErrored = "errored"
	ChJoined  = "joined"
	ChJoing   = "joining"
)

type Chan struct {
	conn   *Conn
	topic  string
	status string
}

func (conn *Conn) Chan(topic string) (*Chan, error) {
	if conn.status != ConnOpen {
		return nil, errors.New("Connection is " + conn.status)
	}

	ch := &Chan{
		conn:   conn,
		topic:  topic,
		status: ChJoing,
	}

	return ch, nil
}

func (ch *Chan) Recv() *MsgCh {
	return ch.conn.register(ch.topic)
}

const (
	ChanClose = "phx_close"
	ChanError = "phx_error"

	// server reply event
	ChanReply = "phx_reply"

	// client send event
	ChanJoin  = "phx_join"
	ChanLeave = "phx_leave"
)

func (ch *Chan) Join() (*MsgCh, error) {
	return ch.Request(ChanJoin, "")
}

func (ch *Chan) Leave() (*MsgCh, error) {
	return ch.Request(ChanLeave, "")
}

func (ch *Chan) On(evt string) *MsgCh {
	return ch.conn.register(ch.topic + evt)
}

func (ch *Chan) Request(evt string, payload interface{}) (*MsgCh, error) {
	msg := &Msg{
		Topic:   ch.topic,
		Event:   evt,
		Ref:     ch.conn.makeRef(),
		Payload: payload,
	}
	mch := ch.conn.register(msg.Ref)
	if err := ch.conn.Send(msg); err != nil {
		mch.Close()
		return nil, err
	}
	return mch, nil
}

func (ch *Chan) Push(evt string, payload interface{}) error {
	msg := &Msg{
		Topic:   ch.topic,
		Event:   evt,
		Ref:     ch.conn.makeRef(),
		Payload: payload,
	}
	return ch.conn.Send(msg)
}
