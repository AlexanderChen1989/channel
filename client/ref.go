package client

import (
	"fmt"
	"sync"
)

type refMaker struct {
	sync.Mutex
	index int
}

func (mk *refMaker) ref() string {
	mk.Lock()
	ref := mk.index
	mk.index++
	mk.Unlock()

	return fmt.Sprint(ref)
}
