package vector

import (
	"unsafe"

	"github.com/koykov/bytealg"
	"github.com/koykov/fastconv"
)

type Flag int

const (
	FlagEscape Flag = iota
)

type Byteptr struct {
	bytealg.Byteptr
	flagEsc bool

	vecPtr uintptr
}

func (p *Byteptr) Bytes() []byte {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		return vec.Helper.ConvertByteptr(p)
	}
	return p.RawBytes()
}

func (p *Byteptr) String() string {
	return fastconv.B2S(p.Bytes())
}

func (p *Byteptr) RawBytes() []byte {
	return p.Byteptr.Bytes()
}

func (p *Byteptr) SetFlag(flag Flag, value bool) {
	switch flag {
	case FlagEscape:
		p.flagEsc = value
	}
}

func (p *Byteptr) GetFlag(flag Flag) bool {
	switch flag {
	case FlagEscape:
		return p.flagEsc
	}
	return false
}

func (p *Byteptr) Reset() {
	p.Byteptr.Reset()
	p.flagEsc = false
	p.vecPtr = 0
}

func (p *Byteptr) indirectVector() *Vector {
	if p.vecPtr == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(p.vecPtr))
}
