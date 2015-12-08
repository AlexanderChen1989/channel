package client

import "sync"

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
	puller := center.pool.getPuller()
	puller.center = center
	puller.key = key
	m[puller] = true
	center.regs[key] = m
	return puller
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

func (center *regCenter) getPullers(key string) []*Puller {
	center.RLock()
	defer center.RUnlock()

	var pullers []*Puller
	for puller, _ := range center.regs[key] {
		pullers = append(pullers, puller)
	}
	return pullers
}
