package client

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConn(t *testing.T) {
	timeout := 500 * time.Millisecond
	conn := NewConn(2, timeout)
	var chs []*MsgCh

	const num = 100
	for i := 0; i < num; i++ {
		chs = append(chs, conn.Recv())
	}

	for i, ch := range chs {
		if i < num/2 {
			ch.Close()
		}
	}

	for i, ch := range chs {
		m := conn.chM[&allMsg]
		if i < num/2 {
			_, ok := m[ch]
			assert.False(t, ok)
		} else {
			_, ok := m[ch]
			assert.True(t, ok)
		}
	}

	var wg sync.WaitGroup
	wg.Add(num)
	for i, ch := range chs {
		go func(i int, ch *MsgCh) {
			defer wg.Done()
			if i < num/2 {
				select {
				case <-ch.ch:
					t.Errorf("Should not receive msg from unregistered MsgCh")
				case <-time.After(timeout + time.Second):
				}
			} else {
				select {
				case <-ch.ch:
				case <-time.After(timeout + time.Second):
					t.Errorf("Should receive msg from registered MsgCh")
				}
			}
		}(i, ch)
	}
	wg.Wait()
}
