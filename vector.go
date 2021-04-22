package vector

import (
	"io"
	"unsafe"
)

type Vector struct {
	buf, src []byte
	bufSS    []string
	addr     uint64
	selfPtr  uintptr
	nodes    []Node
	nodeL    int
	errOff   int
	index    index
}

func (vec *Vector) Parse(_ []byte) error {
	return ErrNotImplement
}

func (vec *Vector) ParseStr(_ string) error {
	return ErrNotImplement
}

func (vec *Vector) ParseCopy(_ []byte) error {
	return ErrNotImplement
}

func (vec *Vector) ParseCopyStr(_ string) error {
	return ErrNotImplement
}

func (vec *Vector) Beautify(_ io.Writer) error {
	return ErrNotImplement
}

func (vec *Vector) Len() int {
	return vec.nodeL
}

func (vec *Vector) ErrorOffset() int {
	return vec.errOff
}

func (vec *Vector) Root() *Node {
	return vec.Get()
}

func (vec *Vector) Exists(key string) bool {
	n := vec.Root()
	return n.Exists(key)
}

func (vec *Vector) AcquireNode(depth int) (r *Node) {
	if vec.nodeL < len(vec.nodes) {
		r = &vec.nodes[vec.nodeL]
		r.Reset()
		vec.nodeL++
	} else {
		r = &Node{typ: TypeUnk}
		vec.nodes = append(vec.nodes, *r)
		vec.nodeL++
	}
	r.vecPtr, r.depth = vec.ptr(), depth
	return
}

func (vec *Vector) AcquireNodeWT(depth int, typ Type) *Node {
	n := vec.AcquireNode(depth)
	n.typ = typ
	return n
}

func (vec *Vector) Reset() {
	if vec.nodeL == 0 {
		return
	}
	_ = vec.nodes[vec.nodeL-1]
	for i := 0; i < vec.nodeL; i++ {
		vec.nodes[i].vecPtr = 0
	}
	vec.buf, vec.src = vec.buf[:0], nil
	vec.bufSS = vec.bufSS[:0]
	vec.addr, vec.nodeL, vec.errOff = 0, 0, 0
	vec.index.reset()
}

func (vec *Vector) ptr() uintptr {
	if vec.selfPtr == 0 {
		vec.selfPtr = uintptr(unsafe.Pointer(vec))
	}
	return vec.selfPtr
}
