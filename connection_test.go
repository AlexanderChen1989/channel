package channel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnection(t *testing.T) {
	conn, err := ConnectTo("http://localhost:4000/socket", nil)
	assert.Nil(t, err)
	defer func() {
		assert.Nil(t, conn.Close())
	}()
	ch, err := conn.JoinTo("rooms:lobby")
	assert.Nil(t, err)
	defer ch.Close()
}
