package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"fmt"
	"math/big"
	"warson-blockchain/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func (p PrivateKey) Sign(data []byte) (*Signature, error) {
	// 生成签名结果 r 和 s
	r, s, err := ecdsa.Sign(rand.Reader, p.key, data)
	if err != nil {
		return nil, err
	}
	return &Signature{
		R: r,
		S: s,
	}, nil
}

func GeneratePrivateKey() PrivateKey {
	// 使用 ecdsa.GenerateKey 函数生成椭圆曲线 P256 上的私钥，随机数由 rand.Reader 提供
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

func (p PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &p.key.PublicKey,
	}
}

type PublicKey struct {
	Key *ecdsa.PublicKey
}

// 实现encoding/gob.GobEncoder接口和gob.GobDecoder，先将PublicKey 转换为字节数组
// 避免gob直接序列化 ecdsa.PublicKey，出现"gob: type elliptic.p256Curve has no exported fields"的错误
// PublicKey 序列化为字节
func (pk *PublicKey) GobEncode() ([]byte, error) {
	// todo:密钥为空，那说明不需要序列化密钥
	if pk.Key == nil {
		return nil, nil
	}
	pubKeyBytes, err := x509.MarshalPKIXPublicKey(pk.Key)
	if err != nil {
		return nil, err
	}
	return pubKeyBytes, nil
}

// PublicKey 反序列化
func (pk *PublicKey) GobDecode(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	pubKey, err := x509.ParsePKIXPublicKey(data)
	if err != nil {
		return err
	}
	ecdsaPubKey, ok := pubKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an ecdsa public key")
	}
	pk.Key = ecdsaPubKey
	return nil
}

func (p PublicKey) ToSlice() []byte {
	// 在椭圆曲线加密中，公钥 PublicKey 是一个点，表示为 (X, Y) 坐标对。
	// 这些坐标是大整数（big.Int 类型），位于曲线 Curve
	return elliptic.MarshalCompressed(p.Key, p.Key.X, p.Key.Y)
}

func (p PublicKey) Address() types.Address {
	h := sha256.Sum256(p.ToSlice())

	// 取后 20 字节的原因：地址更短，更便于存储和传输
	// 以太坊地址就是由公钥的哈希生成的一个 20 字节地址
	// 公钥通常很长，而地址短小精悍，是公钥的“简化版”，便于作为账户或用户的唯一标识
	// 地址隐藏了原始公钥的具体内容，提供了一定程度的隐私
	return types.AddressFromBytes(h[len(h)-20:])
}

type Signature struct {
	R, S *big.Int
}

func (s Signature) Verify(publicKey PublicKey, data []byte) bool {
	// fmt.Printf("Data length: %d and signature is %+v\n", len(data), s)
	return ecdsa.Verify(publicKey.Key, data, s.R, s.S)
}
