package main

import (
	"fmt"
	"log"
	"time"
	"warson-blockchain/crypto"
	"warson-blockchain/network"
)

var transports = []network.Transport{
	network.NewLocalTransport("LOCAL"),
	// network.NewLocalTransport("REMOTE_B"),
	// network.NewLocalTransport("REMOTE_C"),
}

func main() {
	initRemoteServers(transports)
	localNode := transports[0]
	trLate := network.NewLocalTransport("LATE_NODE")

	go func() {
		time.Sleep(7 * time.Second)
		lateServer := makeServer(string(trLate.Addr()), trLate, nil)
		go lateServer.Start()
	}()

	/*
		if err := sendGetStatusMessage(trRemoteA, "REMOTE_B"); err != nil {
			log.Fatal(err)
		}
	*/

	/*
		go func() {
			time.Sleep(7 * time.Second)

			trLate := network.NewLocalTransport("LATE_REMOTE")
			trRemoteC.Connect(trLate)
			lateServer := makeServer(string(trLate.Addr()), trLate, nil)

			go lateServer.Start()
		}()
	*/

	privKey := crypto.GeneratePrivateKey()
	// todo:在initRemoteServer时，已经localServer已经start过了，这里又start一次，会有问题！
	localServer := makeServer("LOCAL", localNode, &privKey)
	localServer.Start()
}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("ID_%s", trs[i].Addr())
		s := makeServer(id, trs[i], nil)
		go s.Start()
	}
}

func makeServer(id string, tr network.Transport, pk *crypto.PrivateKey) *network.Server {
	opts := network.ServerOpts{
		Transport:  tr, // todo:意义何在?代表本server的地址?
		PrivateKey: pk,
		ID:         id,
		Transports: transports, // todo:这里是记录该server知道的其它server
	}

	s, err := network.NewServer(opts)
	if err != nil {
		log.Fatal(err)
	}

	return s
}

/*
func sendTransaction(tr network.Transport, to network.NetAddr) error {
	privKey := crypto.GeneratePrivateKey()
	// todo:直接输入特定指令
	// contract := []byte{0x02, 0x0a, 0x02, 0x0a, 0x0b}
	contract := []byte{0x03, 0x0a, 0x46, 0x0c, 0x4f, 0x0c, 0x4f, 0x0c, 0x0d, 0x05, 0x0a, 0x0f}
	tx := core.NewTransaction(contract)
	tx.Sign(privKey)
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	return tr.SendMessage(to, msg.Bytes())
}

func sendGetStatusMessage(tr network.Transport, to network.NetAddr) error {
	var (
		getStatusMsg = new(network.GetStatusMessage)
		buf          = new(bytes.Buffer)
	)
	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())
	return tr.SendMessage(to, msg.Bytes())
}

*/
