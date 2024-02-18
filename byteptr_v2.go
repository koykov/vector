package vector

import (
	"reflect"
	"unsafe"

	"github.com/koykov/bitset"
	"github.com/koykov/fastconv"
	"github.com/koykov/indirect"
)

type ByteptrV2 struct {
	bits       bitset.Bitset8
	addr, vptr uintptr

	offset, len, cap int
}

func (p *ByteptrV2) TakeAddr(s []byte) *ByteptrV2 {
	if s == nil {
		return p
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&s))
	p.addr, p.cap = h.Data, h.Cap
	return p
}

func (p *ByteptrV2) TakeStringAddr(s string) *ByteptrV2 {
	if len(s) == 0 {
		return p
	}
	h := (*reflect.StringHeader)(unsafe.Pointer(&s))
	p.addr, p.cap = h.Data, h.Len
	return p
}

// DEPRECATED: use TakeStringAddr instead.
func (p *ByteptrV2) TakeStrAddr(s string) *ByteptrV2 {
	return p.TakeStringAddr(s)
}

func (p *ByteptrV2) Bytes() []byte {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		return vec.Helper.Indirect(p)
	}
	return p.RawBytes()
}

func (p *ByteptrV2) String() string {
	if vec := p.indirectVector(); vec != nil && vec.Helper != nil {
		b := vec.Helper.Indirect(p)
		return fastconv.B2S(b)
	}
	return p.RawString()
}

func (p *ByteptrV2) RawBytes() []byte {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return nil
	}
	h := reflect.SliceHeader{
		Data: p.addr + uintptr(p.offset),
		Len:  p.len,
		Cap:  p.len,
	}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func (p *ByteptrV2) RawString() string {
	if p.addr == 0 || p.offset < 0 || p.len < 0 {
		return ""
	}
	h := reflect.StringHeader{
		Data: p.addr + uintptr(p.offset),
		Len:  p.len,
	}
	return *(*string)(unsafe.Pointer(&h))
}

func (p *ByteptrV2) SetBit(pos int, value bool) { p.bits.SetBit(pos, value) }
func (p *ByteptrV2) CheckBit(pos int) bool      { return p.bits.CheckBit(pos) }

func (p *ByteptrV2) Reset() {
	p.bits.Reset()
	p.addr, p.vptr = 0, 0
	p.offset, p.len, p.cap = 0, 0, 0
}

// Restore the entire object from the unsafe pointer.
//
// This needs to reduce pointers count and avoids redundant GC checks.
func (p *ByteptrV2) indirectVector() *Vector {
	if p.vptr == 0 {
		return nil
	}
	return (*Vector)(indirect.ToUnsafePtr(p.vptr))
}
