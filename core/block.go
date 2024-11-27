package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"io"
	"warson-blockchain/crypto"
	"warson-blockchain/types"
)

type Header struct {
	Version       uint32
	DataHash      types.Hash
	PrevBlockHash types.Hash
	Timestamp     uint64
	Height        uint32
}

type Block struct {
	*Header
	Transactions []Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature

	// cached version of header hash
	hash types.Hash
}

func NewBlock(h *Header, txx []Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: txx,
	}
}

func (b *Block) Sign(privateKey crypto.PrivateKey) error {
	fmt.Printf("Data length in Sign: %d\n", len(b.HeaderData())) // 确认签名数据的长度
	sig, err := privateKey.Sign(b.HeaderData())
	if err != nil {
		return err
	}
	b.Validator = privateKey.PublicKey()
	b.Signature = sig

	return nil
}

func (b *Block) Verify() error {
	if b.Signature == nil {
		return fmt.Errorf("block has no signature")
	}

	if !b.Signature.Verify(b.Validator, b.HeaderData()) {
		return fmt.Errorf("block has invalid signature")
	}
	return nil
}

func (b *Block) Encode(r io.Writer, enc Encoder[*Block]) error {
	return enc.Encode(r, b)
}

func (b *Block) Decode(w io.Reader, dec Decoder[*Block]) error {
	return dec.Decode(w, b)
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsZero() {
		b.hash = hasher.Hash(b)
	}
	return b.hash
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	if err := enc.Encode(b.Header); err != nil {
		panic(fmt.Sprintf("Failed to encode header: %v", err))
	}

	// 对编码后的数据进行哈希处理
	hash := sha256.Sum256(buf.Bytes())

	// 返回哈希值
	return hash[:]
}
