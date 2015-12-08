package client

import "sync"

type regCenter struct {
	sync.RWMutex
	pool *pool
	regs map[string][]*Puller
}

func newRegCenter() *regCenter {
	return &regCenter{
		pool: newPool(),
		regs: make(map[string][]*Puller),
	}
}

func (center *regCenter) register(key string) *Puller {
	center.Lock()
	defer center.Unlock()

	puller := center.pool.getPuller()
	puller.center = center
	puller.key = key

	center.regs[key] = append(center.regs[key], puller)

	return puller
}

func (center *regCenter) unregister(puller *Puller) {
	center.Lock()
	defer center.Unlock()

	if puller == nil {
		return
	}

	pullers := center.regs[puller.key]
	for i, _puller := range pullers {
		if _puller != puller {
			continue
		}
		pullers[i] = pullers[len(pullers)-1]
		center.regs[puller.key] = pullers[:len(pullers)-1]
		center.pool.putPuller(puller)
		return
	}

}

func (center *regCenter) getPullers(key string) []*Puller {
	center.RLock()
	defer center.RUnlock()

	pullers := center.regs[key]
	copied := make([]*Puller, len(pullers))
	copy(copied, pullers)

	return copied
}
