package client

import (
	"fmt"
	"sync"
)

type regCenter struct {
	sync.RWMutex
	regs map[string][]*Puller
}

func toKey(topic, event, ref string) string {
	return fmt.Sprintf("KEY:%s:%s:%s:", topic, event, ref)
}

func newRegCenter() *regCenter {
	return &regCenter{
		regs: make(map[string][]*Puller),
	}
}

func (center *regCenter) register(key string) *Puller {
	center.Lock()
	defer center.Unlock()

	puller := &Puller{
		center: center,
		key:    key,
		ch:     make(chan *Message, 1),
	}
	center.regs[key] = append(center.regs[key], puller)

	return puller
}

func (center *regCenter) unregister(puller *Puller) {
	center.Lock()
	defer center.Unlock()

	pullers := center.regs[puller.key]
	for i, _puller := range pullers {
		if _puller != puller {
			continue
		}
		pullers[i] = pullers[len(pullers)-1]
		center.regs[puller.key] = pullers[:len(pullers)-1]
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
