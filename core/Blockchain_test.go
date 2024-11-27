package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func NewBlockChainWithGenesis(t *testing.T) *BlockChain {
	bc, err := NewBlockChain(RandomBlock(0))
	assert.Nil(t, err)
	return bc
}

func TestBlockChain(t *testing.T) {
	bc := NewBlockChainWithGenesis(t)
	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestAddBlock(t *testing.T) {
	bc := NewBlockChainWithGenesis(t)

	lenBlocks := 200
	for i := 0; i < lenBlocks; i++ {
		block := RandomBlockWithSignature(t, uint32(i+1))
		assert.Nil(t, bc.AddBlock(block))
	}

	assert.Equal(t, bc.Height(), uint32(lenBlocks))
	assert.Equal(t, len(bc.headers), lenBlocks+1)
	assert.NotNil(t, bc.AddBlock(RandomBlock(33)))
}
