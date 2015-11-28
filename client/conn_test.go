package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConn(t *testing.T) {
	conn, err := Connect("http://localhost:4000/socket", nil)
	assert.Nil(t, err)
	assert.Nil(t, conn.Close())
}
