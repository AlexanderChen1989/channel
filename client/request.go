package client

type Request struct {
	Parent
	Ref   string
	Msg   *Msg
	msgCh chan *Msg
}

func NewRequest(parent Parent, msg *Msg) *Request {
	return &Request{Parent: parent, Msg: msg}
}

func (req *Request) Reply() chan *Msg {
	req.msgCh = req.Get()
	req.Ref = req.MakeRef()
	req.Pull(req.Ref, req.msgCh)
	req.Msg.Ref = req.Ref
	req.Push(req.Msg)
	return req.msgCh
}

func (req *Request) Done() {
	req.Remove(req.Ref)

	if req.msgCh == nil {
		return
	}

	req.Put(req.msgCh)
}
