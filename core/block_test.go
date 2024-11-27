package core

import (
	"testing"
	"time"
	"warson-blockchain/crypto"
	"warson-blockchain/types"

	"github.com/stretchr/testify/assert"
)

func RandomBlock(height uint32) *Block {
	header := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Height:        height,
		Timestamp:     uint64(time.Now().UnixNano()),
	}

	tx := Transaction{
		Data: []byte("hhh"),
	}

	return NewBlock(header, []Transaction{tx})
}

func TestSignBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := RandomBlock(3)

	assert.Nil(t, b.Sign(privateKey))
	assert.NotNil(t, b.Signature)
}

func TestVerifyBlock(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	b := RandomBlock(5)

	assert.Nil(t, b.Sign(privateKey))
	assert.Nil(t, b.Verify())
	b.Height = 89797
	assert.NotNil(t, b.Verify())

	//otherPrivateKey := crypto.GeneratePrivateKey()
	//b.Validator = otherPrivateKey.PublicKey()
	//assert.NotNil(t, b.Verify())

}
