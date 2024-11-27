package core

import "fmt"

type Validator interface {
	ValidaBlock(*Block) error
}

type BlockValidator struct {
	bc *BlockChain
}

func NewBlockValidator(bc *BlockChain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidaBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already has block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}
	err := b.Verify()
	if err != nil {
		return err
	}

	return nil
}
