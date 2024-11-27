package core

import (
	"testing"
	"warson-blockchain/crypto"

	"github.com/stretchr/testify/assert"
)

func TestSignTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.NotNil(t, tx.Signature)
}

func TestVerifyTransaction(t *testing.T) {
	privateKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}

	assert.Nil(t, tx.Sign(privateKey))
	assert.Nil(t, tx.Verify())
	tx.Data = []byte("fejadjijf")
	assert.NotNil(t, tx.Verify())
	tx.Data = []byte("foo")
	assert.Nil(t, tx.Verify())

	otherPrivateKey := crypto.GeneratePrivateKey()
	tx.PublicKey = otherPrivateKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}
