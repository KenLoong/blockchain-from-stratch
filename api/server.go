package api

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"strconv"
	"warson-blockchain/core"
	"warson-blockchain/types"

	"github.com/go-kit/log"
	"github.com/labstack/echo/v4"
)

type TxResponse struct {
	TxCount uint     `json:"tx_count"`
	Hashes  []string `json:"hashes"`
}

type APIError struct {
	Error string `json:"error"`
}

type Block struct {
	Hash          string     `json:"hash"`
	Version       uint32     `json:"version"`
	DataHash      string     `json:"data_hash"`
	PrevBlockHash string     `json:"prev_block_hash"`
	Height        uint32     `json:"height"`
	Timestamp     int64      `json:"timestamp"`
	Validator     string     `json:"validator"`
	Signature     string     `json:"signature"`
	TxResponse    TxResponse `json:"tx_response"`
}

type ServerConfig struct {
	Logger     log.Logger
	ListenAddr string
}

type Server struct {
	ServerConfig
	bc     *core.Blockchain
	txChan chan *core.Transaction
}

func NewServer(cfg ServerConfig, bc *core.Blockchain, txChan chan *core.Transaction) *Server {
	return &Server{
		ServerConfig: cfg,
		bc:           bc,
		txChan:       txChan,
	}
}

func (s *Server) Start() error {
	e := echo.New()

	e.GET("/block/:hashorid", s.handleGetBlock)
	e.GET("/tx/:hash", s.handleGetTx)
	e.POST("/tx", s.handlePostTx)

	return e.Start(s.ListenAddr)
}

func (s *Server) handleGetTx(c echo.Context) error {
	hash := c.Param("hash")

	b, err := hex.DecodeString(hash)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	tx, err := s.bc.GetTxByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, tx)
}

func (s *Server) handleGetBlock(c echo.Context) error {
	hashOrID := c.Param("hashorid")

	height, err := strconv.Atoi(hashOrID)
	if err == nil { // Treat input as block height
		block, err := s.bc.GetBlock(uint32(height))
		if err != nil {
			return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, intoJSONBlock(block))
	}

	// Otherwise, treat input as block hash
	b, err := hex.DecodeString(hashOrID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	block, err := s.bc.GetBlockByHash(types.HashFromBytes(b))
	if err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	return c.JSON(http.StatusOK, intoJSONBlock(block))
}

func intoJSONBlock(block *core.Block) Block {
	txResponse := TxResponse{
		TxCount: uint(len(block.Transactions)),
		Hashes:  make([]string, len(block.Transactions)),
	}

	for i := 0; i < int(txResponse.TxCount); i++ {
		txResponse.Hashes[i] = block.Transactions[i].Hash(core.TxHasher{}).String()
	}

	return Block{
		Hash:          block.Hash(core.BlockHasher{}).String(),
		Version:       block.Header.Version,
		Height:        block.Header.Height,
		DataHash:      block.Header.DataHash.String(),
		PrevBlockHash: block.Header.PrevBlockHash.String(),
		Timestamp:     block.Header.Timestamp,
		Validator:     block.Validator.Address().String(),
		Signature:     block.Signature.String(),
		TxResponse:    txResponse,
	}
}

func (s *Server) handlePostTx(c echo.Context) error {
	tx := &core.Transaction{}

	// Use JSON decoder instead of gob
	if err := json.NewDecoder(c.Request().Body).Decode(tx); err != nil {
		return c.JSON(http.StatusBadRequest, APIError{Error: err.Error()})
	}

	s.txChan <- tx
	return c.JSON(http.StatusOK, map[string]string{"status": "transaction accepted"})
}
