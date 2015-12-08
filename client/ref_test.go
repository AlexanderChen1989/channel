package client

import (
	"sync"
	"sync/atomic"
	"testing"
)

func BenchmarkAtomic(b *testing.B) {
	var index uint64
	for i := 0; i < b.N; i++ {
		atomic.AddUint64(&index, 1)
	}
}

func BenchmarkMutex(b *testing.B) {
	var index uint64
	var lock sync.Mutex
	for i := 0; i < b.N; i++ {
		lock.Lock()
		index++
		lock.Unlock()
	}
}
