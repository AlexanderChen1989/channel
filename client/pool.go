package client

import "sync"

type pool struct {
	puller *sync.Pool
}

func newPool() *pool {
	return &pool{
		puller: &sync.Pool{
			New: func() interface{} {
				return &Puller{
					ch: make(chan *Message, 1),
				}
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
