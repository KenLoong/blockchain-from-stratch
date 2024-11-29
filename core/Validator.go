package core

import "fmt"

type Validator interface {
	ValidaBlock(*Block) error
}

type BlockValidator struct {
	bc *Blockchain
}

func NewBlockValidator(bc *Blockchain) *BlockValidator {
	return &BlockValidator{
		bc: bc,
	}
}

func (v *BlockValidator) ValidaBlock(b *Block) error {
	if v.bc.HasBlock(b.Height) {
		return fmt.Errorf("chain already has block (%d) with hash (%s)", b.Height, b.Hash(BlockHasher{}))
	}
	if b.Height != v.bc.Height()+1 {
		return fmt.Errorf("block (%s) to height", b.Hash(BlockHasher{}))
	}

	preHeader, err := v.bc.GetHeader(b.Height - 1)
	if err != nil {
		return err
	}

	hash := BlockHasher{}.Hash(preHeader)
	if hash != b.PrevBlockHash {
		return fmt.Errorf("the hash of the previous block (%s) is invalid", b.PrevBlockHash)
	}

	err = b.Verify()
	if err != nil {
		return err
	}

	return nil
}
