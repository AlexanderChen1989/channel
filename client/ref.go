package client

import (
	"fmt"
	"sync/atomic"
)

type refMaker struct {
	index uint64
}

func (mk *refMaker) makeRef() string {
	return fmt.Sprint(atomic.AddUint64(&mk.index, 1))
}
