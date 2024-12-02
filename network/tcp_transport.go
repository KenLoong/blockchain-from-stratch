package network

import (
	"bytes"
	"fmt"
	"net"
)

type TCPPeer struct {
	conn     net.Conn
	Outgoing bool
	ConnAddr string
}

func (p *TCPPeer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	fmt.Printf("conn %s has been sent data\n", p.conn.LocalAddr())
	return err
}

func (p *TCPPeer) readLoop(rpcCh chan RPC) {
	buf := make([]byte, 2048)
	for {
		fmt.Printf("Conn try to read from rpcCh\n")
		n, err := p.conn.Read(buf)
		if err != nil {
			fmt.Printf("read error: %s", err)
			continue
		}

		fmt.Printf("TCPeer Addr %s start send msg to rpc\n", p.ConnAddr)
		msg := buf[:n]
		rpcCh <- RPC{
			From:    p.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}
	}
}

type TCPTransport struct {
	peerCh     chan *TCPPeer
	listenAddr string
	listener   net.Listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		peerCh:     peerCh,
		listenAddr: addr,
	}
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}

	t.listener = ln

	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) readLoop(peer *TCPPeer) {
	buf := make([]byte, 2048)
	for {
		n, err := peer.conn.Read(buf)
		if err != nil {
			fmt.Printf("read error: %s", err)
			continue
		}

		msg := buf[:n]
		fmt.Println(string(msg))
		// handleMessage => server

	}
}

func (t *TCPTransport) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			fmt.Printf("accept error from %+v\n", conn)
			continue
		}

		peer := &TCPPeer{
			conn:     conn,
			ConnAddr: conn.LocalAddr().String(),
		}

		// 注册自己的地址
		t.peerCh <- peer

		// todo:这个日志只打印了一次
		fmt.Printf("new incoming TCP connection => %+v\n", conn)

	}
}
