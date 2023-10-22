package vector

import (
	"github.com/koykov/bitset"
	"github.com/koykov/byteptr"
	"github.com/koykov/fastconv"
	"github.com/koykov/indirect"
)

// Byteptr represents Vector implementation of byteptr.Byteptr object.
type Byteptr struct {
	byteptr.Byteptr
	bitset.Bitset8
	// Vector raw pointer.
	vecPtr uintptr
}

// Bytes converts byteptr object to bytes slice and implements vector's helper logic over it.
func (p *Byteptr) Bytes() []byte {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		return vec.Helper.Indirect(p)
	}
	return p.RawBytes()
}

// Convert byteptr to string.
func (p *Byteptr) String() string {
	return fastconv.B2S(p.Bytes())
}

// RawBytes converts byteptr to byte slice without any logic.
func (p *Byteptr) RawBytes() []byte {
	return p.Byteptr.Bytes()
}

// Reset byteptr object.
func (p *Byteptr) Reset() {
	p.Byteptr.Reset()
	p.Bitset8.Reset()
	p.vecPtr = 0
}

// Restore the entire object from the unsafe pointer.
//
// This needs to reduce pointers count and avoids redundant GC checks.
func (p *Byteptr) indirectVector() *Vector {
	if p.vecPtr == 0 {
		return nil
	}
	return (*Vector)(indirect.ToUnsafePtr(p.vecPtr))
}
