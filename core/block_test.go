package core

import (
	"bytes"
	"fmt"
	"testing"
	"time"
	"warson-blockchain/types"

	"github.com/stretchr/testify/assert"
)

func TestHeaderEncodeAndDecode(t *testing.T) {
	h := &Header{
		Version:   1,
		PreBlock:  types.RandomHash(),
		Timestamp: uint64(time.Now().UnixNano()),
		Height:    10,
		Nonce:     983435,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))
	assert.Equal(t, h, hDecode)
}

func TestBlockEncodeAndDecode(t *testing.T) {
	h := &Header{
		Version:   1,
		PreBlock:  types.RandomHash(),
		Timestamp: uint64(time.Now().UnixNano()),
		Height:    10,
		Nonce:     983435,
	}
	b := &Block{
		Header: *h,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))
	assert.Equal(t, b, bDecode)
}

func TestBlockHash(t *testing.T) {

	h := &Header{
		Version:   1,
		PreBlock:  types.RandomHash(),
		Timestamp: uint64(time.Now().UnixNano()),
		Height:    10,
		Nonce:     983435,
	}
	b := &Block{
		Header: *h,
	}
	hash := b.Hash()
	fmt.Println(hash)
	assert.False(t, hash.IsZero())
}
