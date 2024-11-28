package network

import (
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
