package core

import (
	"encoding/gob"
	"io"
)

type Encoder[T any] interface {
	Encode(T) error
}

type Decoder[T any] interface {
	Decode(T) error
}

type GobTxEncoder struct {
	w io.Writer
}

func NewGobTxEncoder(w io.Writer) *GobTxEncoder {
	//gob.Register(elliptic.P256())
	//gob.Register(elliptic.CurveParams{})
	return &GobTxEncoder{
		w: w,
	}
}

func (e *GobTxEncoder) Encode(tx *Transaction) error {
	return gob.NewEncoder(e.w).Encode(tx)
}

type GobTxDecoder struct {
	r io.Reader
}

func NewGobTxDecoder(r io.Reader) *GobTxDecoder {
	// 默认情况下，Gob 只能处理基本数据类型和已导出字段的结构体。
	// 如果序列化和反序列化中涉及接口或复杂类型，需要提前通过 gob.Register 显式注册这些类型
	// gob.Register 注册的作用域是全局的，它在整个程序运行期间有效
	//gob.Register(elliptic.P256())
	//gob.Register(elliptic.CurveParams{})
	return &GobTxDecoder{
		r: r,
	}
}

func (e *GobTxDecoder) Decode(tx *Transaction) error {
	return gob.NewDecoder(e.r).Decode(tx)
}
