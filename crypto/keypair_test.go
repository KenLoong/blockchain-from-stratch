package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignVerify(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()
	msg := []byte("hello from warson")

	sig, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sig.Verify(publicKey, msg))
	assert.False(t, sig.Verify(publicKey, []byte("xxxxxx")))

	otherPrivateKey := GeneratePrivateKey()
	otherPublicKey := otherPrivateKey.PublicKey()
	assert.False(t, sig.Verify(otherPublicKey, msg))
}
