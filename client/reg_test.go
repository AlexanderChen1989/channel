package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegCenter(t *testing.T) {
	c := newRegCenter()
	const num = 10000
	var chs []*Puller
	for i := 0; i < num; i++ {
		chs = append(chs, c.register(fmt.Sprint(i)))
	}
	for _, ch := range chs {
		assert.Equal(t, len(c.getPullers(ch.key)), 1)
	}
	for _, ch := range chs {
		ch.Close()
	}
	for _, ch := range chs {
		assert.Equal(t, len(c.getPullers(ch.key)), 0)
	}
}
