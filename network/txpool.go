package network

import (
	"sort"
	"warson-blockchain/core"
	"warson-blockchain/types"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

// Len implements sort.Interface.
func (t *TxMapSorter) Len() int {
	return len(t.transactions)
}

// Less implements sort.Interface.
func (t *TxMapSorter) Less(i int, j int) bool {
	return t.transactions[i].FirstSeen() < t.transactions[j].FirstSeen()
}

// Swap implements sort.Interface.
func (t *TxMapSorter) Swap(i int, j int) {
	t.transactions[i], t.transactions[j] = t.transactions[j], t.transactions[i]
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))

	i := 0
	for _, v := range txMap {
		txx[i] = v
		i++
	}

	s := &TxMapSorter{txx}
	sort.Sort(s)
	return s
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (p *TxPool) Transactions() []*core.Transaction {
	s := NewTxMapSorter(p.transactions)
	return s.transactions
}

// Add adds an transaction to the pool, the caller is responsible checking if the
// tx already exist.
func (p *TxPool) Add(tx *core.Transaction) error {
	hash := tx.Hash(core.TxHasher{})
	p.transactions[hash] = tx

	return nil
}

func (p *TxPool) Has(hash types.Hash) bool {
	_, ok := p.transactions[hash]
	return ok
}

func (p *TxPool) Len() int {
	return len(p.transactions)
}

func (p *TxPool) Flush() {
	p.transactions = make(map[types.Hash]*core.Transaction)
}
