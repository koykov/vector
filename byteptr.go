package vector

import (
	"reflect"
	"unsafe"

	"github.com/koykov/bitset"
	"github.com/koykov/byteconv"
	"github.com/koykov/entry"
	"github.com/koykov/indirect"
)

const byteptrSize = 32

type Byteptr struct {
	addr, vptr uintptr

	bits bitset.Bitset32

	offset, len, cap uint32
}

func (p *Byteptr) TakeAddr(s []byte) *Byteptr {
	if s == nil {
		return p
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	p.addr, p.cap = h.Data, uint32(h.Cap)
	return p
}

func (p *Byteptr) TakeStringAddr(s string) *Byteptr {
	if len(s) == 0 {
		return p
	}
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.addr, p.cap = h.Data, uint32(h.Len)
	return p
}

// DEPRECATED: use TakeStringAddr instead.
func (p *Byteptr) TakeStrAddr(s string) *Byteptr {
	return p.TakeStringAddr(s)
}

func (p *Byteptr) Init(s []byte, offset, len int) *Byteptr {
	return p.TakeAddr(s).SetOffset(offset).SetLen(len)
}

func (p *Byteptr) InitRaw(addr uintptr, offset, len int) *Byteptr {
	p.addr = addr
	p.offset = uint32(offset)
	p.len = uint32(len)
	return p
}

func (p *Byteptr) InitString(s string, offset, len int) *Byteptr {
	return p.TakeStringAddr(s).SetOffset(offset).SetLen(len)
}

// DEPRECATED: use InitString instead.
func (p *Byteptr) InitStr(s string, offset, len int) *Byteptr {
	return p.InitString(s, offset, len)
}

func (p *Byteptr) SetAddr(addr uintptr, cap int) *Byteptr {
	p.addr, p.cap = addr, uint32(cap)
	return p
}

func (p *Byteptr) SetOffset(offset int) *Byteptr {
	p.offset = uint32(offset)
	return p
}

func (p *Byteptr) SetLen(len int) *Byteptr {
	p.len = uint32(len)
	return p
}

func (p *Byteptr) Offset() int { return int(p.offset) }
func (p *Byteptr) Len() int    { return int(p.len) }

func (p *Byteptr) Bytes() []byte {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		return vec.Helper.Indirect(p)
	}
	return p.RawBytes()
}

func (p *Byteptr) String() string {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		b := vec.Helper.Indirect(p)
		return byteconv.B2S(b)
	}
	return p.RawString()
}

func (p *Byteptr) RawBytes() []byte {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return nil
	}
	h := reflect.SliceHeader{
		Data: p.addr + uintptr(p.offset),
		Len:  int(p.len),
		Cap:  int(p.len),
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *Byteptr) RawString() string {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return ""
	}
	h := reflect.StringHeader{
		Data: p.addr + uintptr(p.offset),
		Len:  int(p.len),
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *Byteptr) SetBit(pos int, value bool) { p.bits.SetBit(pos, value) }
func (p *Byteptr) CheckBit(pos int) bool      { return p.bits.CheckBit(pos) }

func (p *Byteptr) Reset() {
	p.addr = 0
	p.vptr = 0
	p.bits.Reset()
	p.offset = 0
	p.len = 0
	p.cap = 0
}

// Restore the entire object from the unsafe pointer.
//
// This needs to reduce pointers count and avoids redundant GC checks.
func (p *Byteptr) indirectVector() *Vector {
	if p.vptr == 0 {
		return nil
	}
	return (*Vector)(indirect.ToUnsafePtr(p.vptr))
}

// Check if p equals to substring of s described by e.
func (p *Byteptr) equalKE(s string, e entry.Entry64) bool {
	lo, hi := e.Decode()
	ks := s[lo:hi]
	return ks == p.String()
}
