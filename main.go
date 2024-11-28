package main

import (
	"time"
	"warson-blockchain/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	mockRemote := network.NewLocalTransport("Remote")

	trLocal.Connect(mockRemote)
	mockRemote.Connect(trLocal)

	go func() {
		for {
			mockRemote.SendMessage(trLocal.Addr(), []byte("hello from warson"))
			time.Sleep(2 * time.Second)
		}
	}()

	opts := network.ServerOpts{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(opts)
	s.Start()
}
