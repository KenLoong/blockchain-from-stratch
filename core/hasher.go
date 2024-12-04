package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"warson-blockchain/types"
)

type Hasher[T any] interface {
	Hash(T) types.Hash
}

type BlockHasher struct{}

func (BlockHasher) Hash(b *Header) types.Hash {
	h := sha256.Sum256(b.Bytes())
	return types.Hash(h)
}

type TxHasher struct {
}

func (th TxHasher) Hash(tx *Transaction) types.Hash {
	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(tx); err != nil {
		panic(err)
	}
	return types.Hash(sha256.Sum256(buf.Bytes()))
}
