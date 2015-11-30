package client

import "errors"

const (
	// Channel status

	// ChClosed channel closed
	ChClosed = "closed"
	// ChErrored channel errored
	ChErrored = "errored"
	// ChJoined channel joined
	ChJoined = "joined"
	// ChJoing channel joining
	ChJoing = "joining"

	// Channel events

	// ChanClose channel close event
	ChanClose = "phx_close"
	// ChanError channel error event
	ChanError = "phx_error"
	// ChanReply server reply event
	ChanReply = "phx_reply"
	// ChanJoin  client send join event
	ChanJoin = "phx_join"
	// ChanLeave client send leave event
	ChanLeave = "phx_leave"
)

// Chan channel
type Chan struct {
	conn   *Conn
	topic  string
	status string
}

// Chan create a new channel on connection
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

// Recv register a MsgCh to recv all msg on channel
func (ch *Chan) Recv() *MsgCh {
	return ch.conn.register(ch.topic)
}

// Join channel join, return a MsgCh to receive join result
func (ch *Chan) Join() (*MsgCh, error) {
	return ch.Request(ChanJoin, "")
}

// Leave channel leave, return a MsgCh to receive leave result
func (ch *Chan) Leave() (*MsgCh, error) {
	return ch.Request(ChanLeave, "")
}

// On return a MsgCh to receive all msg on some event on this channel
func (ch *Chan) On(evt string) *MsgCh {
	return ch.conn.register(ch.topic + evt)
}

// Request send a msg to channel and return a MsgCh to receive reply
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

// Push send a msg to channel
func (ch *Chan) Push(evt string, payload interface{}) error {
	msg := &Msg{
		Topic:   ch.topic,
		Event:   evt,
		Ref:     ch.conn.makeRef(),
		Payload: payload,
	}
	return ch.conn.Send(msg)
}
