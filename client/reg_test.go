package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegCenter(t *testing.T) {
	c := newRegCenter()
	const num = 100000
	var chs []*Puller
	for i := 0; i < num; i++ {
		chs = append(chs, c.register(fmt.Sprint(i)))
	}
	for _, ch := range chs {
		assert.Equal(t, 1, len(c.getPullers(ch.key)))
	}
	for _, ch := range chs {
		ch.Close()
	}
	for _, ch := range chs {
		assert.Equal(t, 0, len(c.getPullers(ch.key)))
	}
	for i := 0; i < num; i++ {
		chs = append(chs[:], c.register("Hello"))
	}
	assert.Equal(t, num, len(c.getPullers("Hello")))
	for _, ch := range chs {
		ch.Close()
	}
	assert.Equal(t, 0, len(c.getPullers("Hello")))
}
