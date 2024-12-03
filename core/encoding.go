package core

import (
	"encoding/json"
	"io"
)

// Encoder 定义了一个通用的编码器接口
type Encoder[T any] interface {
	Encode(T) error
}

// Decoder 定义了一个通用的解码器接口
type Decoder[T any] interface {
	Decode(T) error
}

// JSONTxEncoder 使用 JSON 编码交易
type JSONTxEncoder struct {
	w io.Writer
}

// NewJSONTxEncoder 创建一个新的 JSON 交易编码器
func NewJSONTxEncoder(w io.Writer) *JSONTxEncoder {
	return &JSONTxEncoder{
		w: w,
	}
}

// Encode 将交易编码为 JSON
func (e *JSONTxEncoder) Encode(tx *Transaction) error {
	return json.NewEncoder(e.w).Encode(tx)
}

// JSONTxDecoder 使用 JSON 解码交易
type JSONTxDecoder struct {
	r io.Reader
}

// NewJSONTxDecoder 创建一个新的 JSON 交易解码器
func NewJSONTxDecoder(r io.Reader) *JSONTxDecoder {
	return &JSONTxDecoder{
		r: r,
	}
}

// Decode 从 JSON 解码交易
func (e *JSONTxDecoder) Decode(tx *Transaction) error {
	return json.NewDecoder(e.r).Decode(tx)
}

// JSONBlockEncoder 使用 JSON 编码区块
type JSONBlockEncoder struct {
	w io.Writer
}

// NewJSONBlockEncoder 创建一个新的 JSON 区块编码器
func NewJSONBlockEncoder(w io.Writer) *JSONBlockEncoder {
	return &JSONBlockEncoder{
		w: w,
	}
}

// Encode 将区块编码为 JSON
func (enc *JSONBlockEncoder) Encode(b *Block) error {
	return json.NewEncoder(enc.w).Encode(b)
}

// JSONBlockDecoder 使用 JSON 解码区块
type JSONBlockDecoder struct {
	r io.Reader
}

// NewJSONBlockDecoder 创建一个新的 JSON 区块解码器
func NewJSONBlockDecoder(r io.Reader) *JSONBlockDecoder {
	return &JSONBlockDecoder{
		r: r,
	}
}

// Decode 从 JSON 解码区块
func (dec *JSONBlockDecoder) Decode(b *Block) error {
	return json.NewDecoder(dec.r).Decode(b)
}
