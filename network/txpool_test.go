package network

import (
	"math/rand"
	"strconv"
	"testing"
	"warson-blockchain/core"

	"github.com/stretchr/testify/assert"
)

func TestTxPool(t *testing.T) {
	p := NewTxPool()
	assert.Equal(t, 0, p.Len())
}

func TestTxpoolAddtx(t *testing.T) {
	p := NewTxPool()
	tx := core.NewTransaction([]byte("dasd"))
	assert.Nil(t, p.Add(tx))
	assert.Equal(t, 1, p.Len())

	_ = core.NewTransaction([]byte("faidsjaifh"))
	assert.Equal(t, 1, p.Len())
	p.Flush()
	assert.Equal(t, 0, p.Len())

}

func TestSortTransactions(t *testing.T) {
	p := NewTxPool()
	txLen := 1000

	for i := 0; i < txLen; i++ {

		tx := core.NewTransaction([]byte(strconv.FormatInt(int64(i), 10)))
		tx.SetFirstSeen(int64(i * rand.Intn(10000)))
		assert.Nil(t, p.Add(tx))
	}

	txx := p.Transactions()
	for i := 0; i < len(txx)-1; i++ {
		assert.True(t, txx[i].FirstSeen() <= txx[i+1].FirstSeen())
	}

}
