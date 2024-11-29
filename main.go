package main

import (
	"bytes"
	"math/rand"
	"strconv"
	"time"
	"warson-blockchain/core"
	"warson-blockchain/crypto"
	"warson-blockchain/network"

	"github.com/sirupsen/logrus"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	mockRemote := network.NewLocalTransport("Remote")

	trLocal.Connect(mockRemote)
	mockRemote.Connect(trLocal)

	go func() {
		for {
			if err := sendTransaction(mockRemote, trLocal.Addr()); err != nil {
				logrus.Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}()

	privateKey := crypto.GeneratePrivateKey()
	opts := network.ServerOpts{
		PrivateKey: &privateKey,
		ID:         "LOCAL",
		Transports: []network.Transport{trLocal},
	}

	s, err := network.NewServer(opts)
	if err != nil {
		logrus.Fatal(err)
		return
	}
	s.Start()
}

func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000000000)), 10))
	tx := core.NewTransaction(data)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}
