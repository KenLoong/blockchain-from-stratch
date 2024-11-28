package core

import (
	"testing"
	"warson-blockchain/types"

	"github.com/stretchr/testify/assert"
)

func newBlockchainWithGenesis(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(randomBlock(0, types.Hash{}))
	assert.Nil(t, err)

	return bc
}

func getPrevBlockHash(t *testing.T, bc *BlockChain, height uint32) types.Hash {
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

	lenBlocks := 200
	for i := 0; i < lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(randomBlock(33, types.Hash{})))
}
func TestTooHeight(t *testing.T) {
	bc := newBlockchainWithGenesis(t)
	assert.NotNil(t, bc.AddBlock(randomBlockWithSignature(t, 4, types.Hash{})))

}

func TestGetHeader(t *testing.T) {
	bc := newBlockchainWithGenesis(t)

	lenBlocks := 200
	for i := 0; i < lenBlocks; i++ {
		block := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.Nil(t, bc.AddBlock(block))
		header, err := bc.GetHeader(block.Height)
		assert.Nil(t, err)
		assert.Equal(t, block.Header, header)
	}

}
