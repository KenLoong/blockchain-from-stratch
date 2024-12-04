package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
	"warson-blockchain/crypto"
	"warson-blockchain/types"
)

type Header struct {
	Version       uint32     `json:"version"`
	DataHash      types.Hash `json:"data_hash"`
	PrevBlockHash types.Hash `json:"prev_block_hash"`
	Timestamp     int64      `json:"timestamp"`
	Height        uint32     `json:"height"`
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	// Use JSON encoder instead of gob
	enc := json.NewEncoder(buf)
	if err := enc.Encode(h); err != nil {
		panic(fmt.Sprintf("Failed to encode header: %v", err))
	}
	// Hash the encoded data
	hash := sha256.Sum256(buf.Bytes())

	// Return the hash as a slice
	return hash[:]
	// Alternatively, if hash is not needed, use:
	// return buf.Bytes()
}

type Block struct {
	*Header
	Transactions []*Transaction    `json:"transactions"`
	Validator    crypto.PublicKey  `json:"validator"`
	Signature    *crypto.Signature `json:"signature"`

	// Cached version of header hash
	hash types.Hash
}

func NewBlock(h *Header, txx []*Transaction) (*Block, error) {
	return &Block{
		Header:       h,
		Transactions: txx,
	}, nil
}

func (b *Block) AddTransaction(tx *Transaction) {
	b.Transactions = append(b.Transactions, tx)
	hash, _ := CalculateDataHash(b.Transactions)
	b.DataHash = hash
}

func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	sig, err := privateKey.Sign(b.Header.Bytes())
	if err != nil {
		return err
	}
	b.Validator = privateKey.PublicKey()
	b.Signature = sig
	return nil
}

func NewBlockFromPrevHeader(prevHeader *Header, txx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(txx)
	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		Height:        prevHeader.Height + 1,
		DataHash:      dataHash,
		PrevBlockHash: BlockHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
	}

	return NewBlock(header, txx)
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.Header.Bytes()) {
		return fmt.Errorf("block has invalid signature")
	}

	for _, tx := range b.Transactions {
		if err := tx.Verify(); err != nil {
			return err
		}
	}

	dataHash, err := CalculateDataHash(b.Transactions)

	if err != nil {
		return err
	}

	if dataHash != b.DataHash {
		return fmt.Errorf("block (%s) has an invalid data hash", b.Hash(BlockHasher{}))
	}
	return nil
}

func (b *Block) Encode(enc Encoder[*Block]) error {
	// If Encoder still needs to work, ensure JSON is used internally
	return enc.Encode(b)
}

func (b *Block) Decode(dec Decoder[*Block]) error {
	// If Decoder still needs to work, ensure JSON is used internally
	return dec.Decode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b.Header)
	}
	return b.hash
}

func CalculateDataHash(txx []*Transaction) (hash types.Hash, err error) {
	buf := &bytes.Buffer{}

	for _, tx := range txx {
		// Use JSON encoder instead of gob
		if err = json.NewEncoder(buf).Encode(tx); err != nil {
			return
		}
	}
	hash = sha256.Sum256(buf.Bytes())
	return
}
