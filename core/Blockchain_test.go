package core

import (
	"fmt"
	"testing"
	"warson-blockchain/crypto"
	"warson-blockchain/types"

	"github.com/go-kit/log"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *Blockchain {
	bc, err := NewBlockchain(log.NewNopLogger(), randomBlock(t, 0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	prevHeader, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)
	return BlockHasher{}.Hash(prevHeader)
}

func TestBlockChain(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestAddBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 10
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 33, types.Hash{})))
}
func TestTooHeight(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.AddBlock(randomBlock(t, 4, types.Hash{})))

}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 200
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, block.Header, header)
	}

}

func TestGetBlock(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	lenBlocks := 100
	for i := 0; i < lenBlocks; i++ {
		block := randomBlock(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		fetchedBlock, err := bc.GetBlock(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, fetchedBlock, block)
	}
}

func TestSendNativeTransferTamper(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))
	privKeyBob := crypto.GeneratePrivateKey()
	privKeyAlice := crypto.GeneratePrivateKey()
	amount := uint64(100)
	accountBob := bc.accountState.CreateAccount(privKeyBob.PublicKey().Address())
	accountBob.Balance = amount
	tx := NewTransaction([]byte{})
	tx.From = privKeyBob.PublicKey()
	tx.To = privKeyAlice.PublicKey()
	tx.Value = amount
	tx.Sign(privKeyBob)
	// 修改sign后的tx
	hackerPrivKey := crypto.GeneratePrivateKey()
	tx.To = hackerPrivKey.PublicKey()
	block.AddTransaction(tx)
	assert.NotNil(t, bc.AddBlock(block)) // this should fail
	//fmt.Printf("%+v\n", hackerPrivKey.PublicKey().Address())
	fmt.Printf("%+v\n", bc.accountState.accounts)
	fmt.Printf("%+v\n", privKeyAlice.PublicKey().Address())
	_, err := bc.accountState.GetAccount(hackerPrivKey.PublicKey().Address())
	assert.Equal(t, ErrAccountNotFound, err)
	_, err = bc.accountState.GetAccount(privKeyAlice.PublicKey().Address())
	assert.NotNil(t, err)
	// assert.Equal(t, accountAlice.Balance, amount)
}
func TestSendNativeTransferInsuffientBalance(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))
	privKeyBob := crypto.GeneratePrivateKey()
	privKeyAlice := crypto.GeneratePrivateKey()
	amount := uint64(100)
	accountBob := bc.accountState.CreateAccount(privKeyBob.PublicKey().Address())
	accountBob.Balance = uint64(99)
	tx := NewTransaction([]byte{})
	tx.From = privKeyBob.PublicKey()
	tx.To = privKeyAlice.PublicKey()
	tx.Value = amount
	tx.Sign(privKeyBob)
	block.AddTransaction(tx)
	assert.NotNil(t, bc.AddBlock(block))
	_, err := bc.accountState.GetAccount(privKeyAlice.PublicKey().Address())
	assert.NotNil(t, err)
}
func TestSendNativeTransferSuccess(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	privKeyBob := crypto.GeneratePrivateKey()
	privKeyAlice := crypto.GeneratePrivateKey()
	amount := uint64(100)
	accountBob := bc.accountState.CreateAccount(privKeyBob.PublicKey().Address())
	accountBob.Balance = amount
	tx := NewTransaction([]byte{})
	tx.From = privKeyBob.PublicKey()
	tx.To = privKeyAlice.PublicKey()
	tx.Value = amount
	assert.Nil(t, tx.Sign(privKeyBob))
	block.AddTransaction(tx)
	assert.Nil(t, block.Sign(signer))
	assert.Nil(t, bc.AddBlock(block))
	accountAlice, err := bc.accountState.GetAccount(privKeyAlice.PublicKey().Address())
	assert.Nil(t, err)
	assert.Equal(t, amount, accountAlice.Balance)
}

func TestBlockSignAndValidate(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	signer := crypto.GeneratePrivateKey()
	block := randomBlock(t, uint32(1), getPrevBlockHash(t, bc, uint32(1)))
	assert.Nil(t, block.Sign(signer))
	assert.Nil(t, bc.AddBlock(block))
}
