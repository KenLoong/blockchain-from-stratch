package network

/*
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
	buf := make([]byte, len(msg))
	n, err := rpc.Payload.Read(buf)
	assert.Nil(t, err)
	assert.Equal(t, len(msg), n)

	assert.Equal(t, msg, buf)
	assert.Equal(t, rpc.From, ta.Addr())

}

func TestBroadcast(t *testing.T) {
	ta := NewLocalTransport("A")
	tb := NewLocalTransport("B")
	tc := NewLocalTransport("C")

	ta.Connect(tb)
	ta.Connect(tc)

	msg := []byte("hello")
	assert.Nil(t, ta.Broadcast(msg))

	rpcB := <-tb.Consume()
	b, err := ioutil.ReadAll(rpcB.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, b)

	rpcC := <-tc.Consume()
	b, err = ioutil.ReadAll(rpcC.Payload)
	assert.Nil(t, err)
	assert.Equal(t, msg, b)
}
*/
