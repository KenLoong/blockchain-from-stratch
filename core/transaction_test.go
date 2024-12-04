package core

import (
	"bytes"
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
	tx.From = otherPrivateKey.PublicKey()

	assert.NotNil(t, tx.Verify())
}

func randomTxWithSignature(t *testing.T) *Transaction {
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		Data: []byte("foo"),
	}
	assert.Nil(t, tx.Sign(privKey))
	return tx
}

func TestTxEncodeDecode(t *testing.T) {
	tx := randomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewJSONTxEncoder(buf)))

	txDecoded := new(Transaction)
	assert.Nil(t, txDecoded.Decode(NewJSONTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func TestNFTTX(t *testing.T) {
	collectionTx := CollectionTx{
		Fee:      200,
		MetaData: []byte("The beginning of a new collection"),
	}
	privKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		TxInner: collectionTx,
	}
	tx.Sign(privKey)
	buf := new(bytes.Buffer)
	assert.Nil(t, tx.Encode(NewJSONTxEncoder(buf)))
	txDecoded := &Transaction{}
	assert.Nil(t, txDecoded.Decode(NewJSONTxDecoder(buf)))
	assert.Equal(t, tx, txDecoded)
}

func TestNativeTransferTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	toPrivKey := crypto.GeneratePrivateKey()
	tx := &Transaction{
		To:    toPrivKey.PublicKey(),
		Value: 666,
	}
	assert.Nil(t, tx.Sign(fromPrivKey))
}
