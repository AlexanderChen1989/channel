package client

import "sync"

type Puller struct {
	center *regCenter
	ch     chan *Msg
	key    string
}

type regCenter struct {
	sync.RWMutex
	pool *pool
	regs map[string]map[*MsgCh]bool
}

func newRegCenter() *regCenter {
	return &regCenter{
		pool: newPool(),
		regs: make(map[string]map[*MsgCh]bool),
	}
}

func (center *regCenter) register(key string) *MsgCh {
	center.Lock()
	defer center.Unlock()

	m := center.regs[key]
	if m == nil {
		m = center.pool.getMsgChMap()
	}
	mch := center.pool.getMsgCh()
	mch.key = key
	m[mch] = true
	center.regs[key] = m
	return mch
}

func (center *regCenter) unregister(key string, mch *MsgCh) {
	center.Lock()
	defer center.Unlock()

	if m := center.regs[key]; m != nil {
		delete(m, mch)
		center.pool.putMsgCh(mch)
		if len(m) == 0 {
			delete(center.regs, key)
			center.pool.putMsgChMap(m)
		}
	}
}

func (center *regCenter) get(key string) map[*MsgCh]bool {
	center.RLock()
	m := center.regs[key]
	center.RUnlock()
	return m
}
