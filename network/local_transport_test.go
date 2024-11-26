package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	ta := NewLocalTransport("A")
	tb := NewLocalTransport("B")

	tb.Connect(ta)
	ta.Connect(tb)
	assert.Equal(t, ta.GetPeer(tb.Addr()), tb)
	assert.Equal(t, tb.GetPeer(ta.Addr()), ta)

}

func TestSendMessage(t *testing.T) {
	ta := NewLocalTransport("A")
	tb := NewLocalTransport("B")

	ta.Connect(tb)
	tb.Connect(ta)

	msg := []byte("hello")
	assert.Nil(t, ta.SendMessage(tb.Addr(), msg))

	rpc := <-tb.Consume()
	assert.Equal(t, rpc.Payload, msg)
	assert.Equal(t, rpc.From, ta.Addr())
}
