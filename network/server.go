package network

import (
	"bytes"
	"encoding/gob"

	"fmt"
	"os"
	"time"
	"warson-blockchain/core"
	"warson-blockchain/crypto"
	"warson-blockchain/types"

	"github.com/go-kit/log"
)

var defaultBlockTime = 5 * time.Second

type ServerOpts struct {
	Transport     Transport
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	ServerOpts
	memPool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcCh       chan RPC
	quitCh      chan struct{}
}

func NewServer(opts ServerOpts) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}
	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		//	opts.Logger = log.With(opts.Logger, "ID", opts.ID)
		opts.Logger = log.With(opts.Logger, "addr", opts.Transport.Addr())
	}

	chain, err := core.NewBlockChain(opts.Logger, genesisBlock())
	if err != nil {
		return nil, err
	}

	s := &Server{
		ServerOpts:  opts,
		memPool:     NewTxPool(1000),
		chain:       chain,
		isValidator: opts.PrivateKey != nil,
		rpcCh:       make(chan RPC),
		quitCh:      make(chan struct{}, 1),
	}

	// If we dont got any processor from the server options, we going to use
	// the server as default.
	if s.RPCProcessor == nil {
		s.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	s.boostrapNodes()

	return s, nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcCh:
			msg, err := s.RPCDecodeFunc(rpc)
			if err != nil {
				s.Logger.Log("error", err)
			}

			if err := s.RPCProcessor.ProcessMessage(msg); err != nil {
				if err != core.ErrBlockKnown {
					s.Logger.Log("error", err)
				}
			}

		case <-s.quitCh:
			break free
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) ProcessMessage(msg *DecodedMessage) error {
	switch t := msg.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(t)
	case *core.Block:
		return s.processBlock(t)
	case *GetStatusMessage:
		return s.processGetStatusMessage(msg.From, t)
	case *StatusMessage:
		return s.processStatusMessage(msg.From, t)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(msg.From, t)
	}
	return nil
}

// TODO(@warson): Remove the logic from the main function to here
// Normally Transport which is our own transport should do the trick.
func (s *Server) sendGetStatusMessage(tr Transport) error {
	var (
		getStatusMsg = new(GetStatusMessage)
		buf          = new(bytes.Buffer)
	)
	if err := gob.NewEncoder(buf).Encode(getStatusMsg); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
	if err := s.Transport.SendMessage(tr.Addr(), msg.Bytes()); err != nil {
		return err
	}

	/// statusMessage
	return nil
}

func (s *Server) broadcast(payload []byte) error {
	for _, tr := range s.Transports {
		if err := tr.Broadcast(payload); err != nil {
			return err
		}
	}
	return nil
}

func (s *Server) processTransaction(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})

	if s.memPool.Contains(hash) {
		return nil
	}

	if err := tx.Verify(); err != nil {
		return err
	}

	/*
		s.Logger.Log(
			"msg", "add new tx to mempool",
			"hash", hash,
			"mempoolLength", s.memPool.Len(),
		)
	*/
	go s.broadcastTx(tx)

	s.memPool.Add(tx)
	return nil
}

// todo: @warson
func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	if err := b.Encode(core.NewGobBlockEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	if err := tx.Encode(core.NewGobTxEncoder(buf)); err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())

	return s.broadcast(msg.Bytes())
}

func (s *Server) createNewBlock() error {
	currentHeader, err := s.chain.GetHeader(s.chain.Height())
	if err != nil {
		return err
	}

	// For now we are going to use all transactions that are in the pending pool
	// Later on when we know the internal structure of our transaction
	// we will implement some kind of complexity function to determine how
	// many transactions can be included in a block.
	txx := s.memPool.Pending()

	block, err := core.NewBlockFromPrevHeader(currentHeader, txx)
	if err != nil {
		return err
	}

	if err := block.Sign(*s.PrivateKey); err != nil {
		return err
	}

	if err := s.chain.AddBlock(block); err != nil {
		return err
	}

	// TODO(@warson): pending pool of tx should only reflect on validator nodes.
	// Right now "normal nodes" does not have their pending pool cleared.
	s.memPool.ClearPending()
	go s.broadcastBlock(block)
	return nil
}

func (s *Server) initTransports() {
	// 这里为什么要监听transports呢？应该监听s.Transport.peers的才对
	// todo:需要修改
	for _, tr := range s.Transports {
		go func(tr Transport) {
			for rpc := range tr.Consume() {
				s.rpcCh <- rpc
			}
		}(tr)
	}
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.BlockTime)
	s.Logger.Log("msg", "Starting validator loop", "blocktime", s.BlockTime)

	for {
		<-ticker.C
		s.createNewBlock()
	}
}

func genesisBlock() *core.Block {
	header := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Height:    0,
		Timestamp: 000000,
	}

	b, _ := core.NewBlock(header, nil)
	return b
}

func (s *Server) processBlock(b *core.Block) error {
	if err := s.chain.AddBlock(b); err != nil {
		return err
	}

	go s.broadcastBlock(b)

	return nil
}

func (s *Server) processStatusMessage(from NetAddr, data *StatusMessage) error {
	if data.CurrentHeight <= s.chain.Height() {
		s.Logger.Log("msg", "cannot sync blockHeight to low", "ourHeight", s.chain.Height(), "theirHeight", data.CurrentHeight, "addr", from)
		return nil
	}

	// In this case we are 100% sure that the node has blocks heigher than us.
	getBlocksMessage := &GetBlocksMessage{
		From: s.chain.Height(),
		To:   0,
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(getBlocksMessage); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
	return s.Transport.SendMessage(from, msg.Bytes())
}

func (s *Server) processGetStatusMessage(from NetAddr, data *GetStatusMessage) error {
	fmt.Printf("=> received Getstatus msg from %s => %+v\n", from, data)
	statusMessage := &StatusMessage{
		CurrentHeight: s.chain.Height(),
		ID:            s.ID,
	}
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(statusMessage); err != nil {
		return err
	}
	msg := NewMessage(MessageTypeStatus, buf.Bytes())
	return s.Transport.SendMessage(from, msg.Bytes())
}

// todo:这个函数很有问题，err不为空的情况下，应该停止接下来的动作才对
func (s *Server) boostrapNodes() {
	for _, tr := range s.Transports {
		if s.Transport.Addr() != tr.Addr() {
			if err := s.Transport.Connect(tr); err != nil {
				s.Logger.Log("error", "could not connect to remote", "err", err)
			}
			s.Logger.Log("msg", "connect to remote", "we", s.Transport.Addr(), "addr", tr.Addr())
			// Send the getStatusMessage so we can sync (if needed)
			if err := s.sendGetStatusMessage(tr); err != nil {
				s.Logger.Log("error", "sendGetStatusMessage", "err", err)
			}
		}
	}
}

func (s *Server) processGetBlocksMessage(from NetAddr, data *GetBlocksMessage) error {
	panic("here")
	fmt.Printf("got get blocks message => %+v\n", data)
	return nil
}
