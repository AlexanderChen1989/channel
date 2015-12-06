package client

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegCenter(t *testing.T) {
	c := newRegCenter()
	const num = 10
	var chs []*MsgCh
	for i := 0; i < num; i++ {
		chs = append(chs, c.register(fmt.Sprint(i)))
	}
	for _, ch := range chs {
		assert.NotNil(t, c.get(ch.key.(string)))
	}
	for _, ch := range chs {
		c.unregister(ch.key.(string), ch)
	}
	for _, ch := range chs {
		assert.Nil(t, c.get(ch.key.(string)))
	}
}
