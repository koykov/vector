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

// Vector implementation of bytealg.Byteptr object.
type Byteptr struct {
	bytealg.Byteptr
	// Escape flag.
	flagEsc bool
	// Vector raw pointer.
	vecPtr uintptr
}

// Convert byteptr object to bytes slice and implements vector's helper logic over it.
func (p *Byteptr) Bytes() []byte {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		return vec.Helper.ConvertByteptr(p)
	}
	return p.RawBytes()
}

// Convert byteptr to string.
func (p *Byteptr) String() string {
	return fastconv.B2S(p.Bytes())
}

// Convert byteptr to byte slice without any logic.
func (p *Byteptr) RawBytes() []byte {
	return p.Byteptr.Bytes()
}

// Set flag.
func (p *Byteptr) SetFlag(flag Flag, value bool) {
	switch flag {
	case FlagEscape:
		p.flagEsc = value
	}
}

// Get flag.
func (p *Byteptr) GetFlag(flag Flag) bool {
	switch flag {
	case FlagEscape:
		return p.flagEsc
	}
	return false
}

// Reset byteptr object.
func (p *Byteptr) Reset() {
	p.Byteptr.Reset()
	p.flagEsc = false
	p.vecPtr = 0
}

// Restore the entire object from the unsafe pointer.
//
// This needs to reduce pointers count and avoids redundant GC checks.
func (p *Byteptr) indirectVector() *Vector {
	if p.vecPtr == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(p.vecPtr))
}
