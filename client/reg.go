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
	regs map[string]map[*Puller]bool
}

func newRegCenter() *regCenter {
	return &regCenter{
		pool: newPool(),
		regs: make(map[string]map[*Puller]bool),
	}
}

func (center *regCenter) register(key string) *Puller {
	center.Lock()
	defer center.Unlock()

	m := center.regs[key]
	if m == nil {
		m = center.pool.getPullerMap()
	}
	mch := center.pool.getPuller()
	mch.key = key
	m[mch] = true
	center.regs[key] = m
	return mch
}

func (center *regCenter) unregister(puller *Puller) {
	center.Lock()
	defer center.Unlock()

	if puller == nil {
		return
	}

	if m := center.regs[puller.key]; m != nil {
		delete(m, puller)
		center.pool.putPuller(puller)
		if len(m) == 0 {
			delete(center.regs, puller.key)
			center.pool.putPullerMap(m)
		}
	}
}

func (center *regCenter) get(key string) map[*Puller]bool {
	center.RLock()
	m := center.regs[key]
	center.RUnlock()
	return m
}
