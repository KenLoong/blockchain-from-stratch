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

func TestPublicKeySerialization(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()

	// 序列化
	encodedData, err := publicKey.GobEncode()
	assert.Nil(t, err)
	// 反序列化
	decodedPubKey := &PublicKey{}
	err = decodedPubKey.GobDecode(encodedData)
	assert.Nil(t, err)

	// 验证解码后是否与原始公钥一致
	assert.Equal(t, publicKey.Key, decodedPubKey.Key)
}
