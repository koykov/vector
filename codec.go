package vector

import (
	"io"
	"strconv"
)

// Codec object interface.
type Codec interface {
	// Decode convert byteptr to byte slice and apply custom logic.
	Decode(p *Byteptr) ([]byte, error)
	// DecodeInt convert byteptr to int value.
	DecodeInt(p *Byteptr) (int64, error)
	// DecodeUint convert byteptr to uint value.
	DecodeUint(p *Byteptr) (uint64, error)
	// DecodeFloat convert byteptr to float value.
	DecodeFloat(p *Byteptr) (float64, error)
	// DecodeBool convert byteptr to bool value.
	DecodeBool(p *Byteptr) (bool, error)

	// Beautify makes a beauty view of node.
	Beautify(w io.Writer, node *Node) error
	// Marshal serializes node.
	Marshal(w io.Writer, node *Node) error
}

type BaseCodec struct{}

func (BaseCodec) Decode(p *Byteptr) ([]byte, error) { return p.RawBytes(), nil }

func (BaseCodec) DecodeInt(p *Byteptr) (int64, error) {
	return strconv.ParseInt(p.RawString(), 10, 64)
}

func (BaseCodec) DecodeUint(p *Byteptr) (uint64, error) {
	return strconv.ParseUint(p.RawString(), 10, 64)
}

func (BaseCodec) DecodeFloat(p *Byteptr) (float64, error) {
	return strconv.ParseFloat(p.RawString(), 64)
}

func (BaseCodec) DecodeBool(p *Byteptr) (bool, error) {
	return strconv.ParseBool(p.RawString())
}

func (BaseCodec) Beautify(_ io.Writer, _ *Node) error {
	return ErrNotImplement
}

func (BaseCodec) Marshal(_ io.Writer, _ *Node) error {
	return ErrNotImplement
}
