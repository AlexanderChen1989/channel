package client

import "sync"

type pool struct {
	puller    *sync.Pool
	pullerMap *sync.Pool
}

func newPool() *pool {
	return &pool{
		puller: &sync.Pool{
			New: func() interface{} {
				return &Puller{
					ch: make(chan *Msg, 1),
				}
			},
		},
		pullerMap: &sync.Pool{
			New: func() interface{} {
				return make(map[*Puller]bool)
			},
		},
	}
}

func (p *pool) getPuller() *Puller {
	return p.puller.Get().(*Puller)
}

func (p *pool) putPuller(puller *Puller) {
	if puller == nil {
		return
	}
	select {
	case <-puller.ch:
	default:
	}
	p.puller.Put(puller)
}

func (p *pool) getPullerMap() map[*Puller]bool {
	return p.pullerMap.Get().(map[*Puller]bool)
}

func (p *pool) putPullerMap(m map[*Puller]bool) {
	if m == nil {
		return
	}
	p.pullerMap.Put(m)
}
