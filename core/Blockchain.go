package core

type BlockChain struct {
	store     Storage
	headers   []*Header
	validator Validator
}

func NewBlockChain(genesis *Block) (*BlockChain, error) {
	bc := &BlockChain{
		headers: []*Header{},
		store:   NewMemoryStorage(),
	}
	bc.validator = NewBlockValidator(bc)
	err := bc.addBlockWithoutValidation(genesis)
	return bc, err
}

func (bc *BlockChain) SetValidator(v Validator) {
	bc.validator = v
}

func (bc *BlockChain) AddBlock(b *Block) error {
	err := bc.validator.ValidaBlock(b)
	if err != nil {
		return err
	}
	return bc.addBlockWithoutValidation(b)
}

func (bc *BlockChain) HasBlock(height uint32) bool {
	return height <= bc.Height()
}

func (bc *BlockChain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *BlockChain) addBlockWithoutValidation(b *Block) error {
	bc.headers = append(bc.headers, b.Header)
	return bc.store.Put(b)
}
