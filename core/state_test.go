package core

import (
	"testing"
	"warson-blockchain/crypto"

	"github.com/stretchr/testify/assert"
)

func TestAccountStateTransferNoBalance(t *testing.T) {
	state := NewAccountState()

	from := crypto.GeneratePrivateKey().PublicKey().Address()
	to := crypto.GeneratePrivateKey().PublicKey().Address()
	amount := uint64(90)

	assert.NotNil(t, state.Transfer(from, to, amount))
}

func TestAccountStateTransferSuccess(t *testing.T) {
	state := NewAccountState()
	from := crypto.GeneratePrivateKey().PublicKey().Address()

	state.CreateAccount(from).Balance = 99
	to := crypto.GeneratePrivateKey().PublicKey().Address()
	amount := uint64(90)

	assert.Nil(t, state.Transfer(from, to, amount))
}
