package client

import "sync"

type pool struct {
	msgCh    *sync.Pool
	msgChMap *sync.Pool
}

func newPool() *pool {
	return &pool{
		msgCh: &sync.Pool{
			New: func() interface{} {
				return &MsgCh{
					ch: make(chan *Msg, 1),
				}
			},
		},
		msgChMap: &sync.Pool{
			New: func() interface{} {
				return make(map[*MsgCh]bool)
			},
		},
	}
}

func (p *pool) getMsgCh() *MsgCh {
	return p.msgCh.Get().(*MsgCh)
}

func (p *pool) putMsgCh(mch *MsgCh) {
	if mch == nil {
		return
	}
	select {
	case <-mch.Recv():
	default:
	}
	p.msgCh.Put(mch)
}

func (p *pool) getMsgChMap() map[*MsgCh]bool {
	return p.msgChMap.Get().(map[*MsgCh]bool)
}

func (p *pool) putMsgChMap(m map[*MsgCh]bool) {
	if m == nil {
		return
	}
	p.msgChMap.Put(m)
}
