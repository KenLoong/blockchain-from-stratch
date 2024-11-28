package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
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
	// 生成密钥对
	privKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	originalPubKey := &privKey.PublicKey
	publicKeyWrapper := &PublicKey{Key: originalPubKey}

	// 序列化
	encodedData, err := publicKeyWrapper.GobEncode()
	if err != nil {
		t.Fatalf("Failed to encode public key: %v", err)
	}

	// 反序列化
	decodedPubKey := &PublicKey{}
	if err := decodedPubKey.GobDecode(encodedData); err != nil {
		t.Fatalf("Failed to decode public key: %v", err)
	}

	// 验证解码后是否与原始公钥一致
	if !originalPubKey.Equal(decodedPubKey.Key) {
		t.Fatalf("Decoded public key does not match the original")
	}
}
