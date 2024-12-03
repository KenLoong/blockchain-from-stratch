package core

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"warson-blockchain/crypto"
	"warson-blockchain/types"
)

type TxType byte

const (
	TxTypeCollection TxType = iota // 0x0
	TxTypeMint                     // 0x01
)

type CollectionTx struct {
	Fee      int64
	MetaData []byte
}

type MintTx struct {
	Fee             int64
	NFT             types.Hash
	Collection      types.Hash
	MetaData        []byte
	CollectionOwner crypto.PublicKey
	Signature       crypto.Signature
}

type Transaction struct {
	Data      []byte
	Type      TxType
	TxInner   any
	From      crypto.PublicKey
	Signature *crypto.Signature
	Nonce     int64
	hash      types.Hash // cached version of the tx data hash
}

// NewTransaction creates a new transaction with random nonce
func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data:  data,
		Nonce: rand.Int63n(10000000000000),
	}
}

// MarshalJSON provides custom encoding for Transaction
func (tx *Transaction) MarshalJSON() ([]byte, error) {
	// 强制指定
	if tx.Type != TxTypeCollection && tx.Type != TxTypeMint {
		tx.Type = TxTypeCollection
	}
	// 如果 TxInner 为 nil，尝试默认填充 TxInner
	if tx.TxInner == nil {
		if tx.Type == TxTypeCollection {
			tx.TxInner = CollectionTx{} // 默认填充 TxTypeCollection
		} else if tx.Type == TxTypeMint {
			tx.TxInner = MintTx{} // 默认填充 TxTypeMint
		}
	}

	type Alias Transaction
	aux := &struct {
		TxInner json.RawMessage `json:"tx_inner"`
		*Alias
	}{
		Alias: (*Alias)(tx),
	}

	// 根据 TxType 编码 TxInner 字段
	switch tx.Type {
	case TxTypeCollection:
		// 如果是 CollectionTx 类型，进行自定义编码
		if collectionTx, ok := tx.TxInner.(CollectionTx); ok {
			txInner, err := json.Marshal(collectionTx)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal CollectionTx: %w", err)
			}
			aux.TxInner = txInner
		}
	case TxTypeMint:
		// 如果是 MintTx 类型，进行自定义编码
		if mintTx, ok := tx.TxInner.(MintTx); ok {
			txInner, err := json.Marshal(mintTx)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal MintTx: %w", err)
			}
			aux.TxInner = txInner
		}
	}

	return json.Marshal(aux)
}

// UnmarshalJSON provides custom decoding for Transaction
func (tx *Transaction) UnmarshalJSON(data []byte) error {
	type Alias Transaction
	aux := &struct {
		TxInner json.RawMessage `json:"tx_inner"`
		*Alias
	}{
		Alias: (*Alias)(tx),
	}

	// 先解码常规字段
	if err := json.Unmarshal(data, aux); err != nil {
		return err
	}

	// 如果 TxInner 是 nil，则需要根据 TxType 来填充默认的 TxInner
	if tx.TxInner == nil {
		if tx.Type == TxTypeCollection {
			tx.TxInner = CollectionTx{} // 默认填充 TxTypeCollection
		} else if tx.Type == TxTypeMint {
			tx.TxInner = MintTx{} // 默认填充 TxTypeMint
		}
	}

	// 根据 TxType 解码 TxInner 字段
	switch tx.Type {
	case TxTypeCollection:
		var collectionTx CollectionTx
		if err := json.Unmarshal(aux.TxInner, &collectionTx); err != nil {
			return fmt.Errorf("failed to unmarshal CollectionTx: %w", err)
		}
		tx.TxInner = collectionTx
	case TxTypeMint:
		var mintTx MintTx
		if err := json.Unmarshal(aux.TxInner, &mintTx); err != nil {
			return fmt.Errorf("failed to unmarshal MintTx: %w", err)
		}
		tx.TxInner = mintTx
	default:
		return fmt.Errorf("unsupported TxType: %v", tx.Type)
	}

	return nil
}

// Sign signs the transaction with the private key
func (tx *Transaction) Sign(privateKey crypto.PrivateKey) error {
	sig, err := privateKey.Sign(tx.Data)
	if err != nil {
		return err
	}

	tx.From = privateKey.PublicKey()
	tx.Signature = sig

	return nil
}

// Verify verifies the transaction signature
func (tx *Transaction) Verify() error {
	if tx.Signature == nil {
		return fmt.Errorf("transaction has no signature")
	}

	if !tx.Signature.Verify(tx.From, tx.Data) {
		return fmt.Errorf("invalid transaction signature")
	}
	return nil
}

func (tx *Transaction) Decode(dec Decoder[*Transaction]) error {
	return dec.Decode(tx)
}

func (tx *Transaction) Encode(dec Encoder[*Transaction]) error {
	return dec.Encode(tx)
}

// Hash calculates the hash of the transaction
func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsZero() {
		tx.hash = hasher.Hash(tx)
	}
	return tx.hash
}
