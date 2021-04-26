package vector

import (
	"io"
	"reflect"
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
	Index    Index
	Helper   Helper
}

type Helper interface {
	ConvertByteptr(*Byteptr) []byte
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

func (vec *Vector) SetSrc(s []byte, copy bool) error {
	if len(s) == 0 {
		return ErrEmptySrc
	}
	if copy {
		vec.buf = append(vec.buf[:0], s...)
		vec.src = vec.buf
	} else {
		vec.src = s
	}
	h := (*reflect.SliceHeader)(unsafe.Pointer(&vec.src))
	vec.addr = uint64(h.Data)
	return nil
}

func (vec *Vector) Len() int {
	return vec.nodeL
}

func (vec *Vector) SrcLen() int {
	return len(vec.src)
}

func (vec *Vector) Src() []byte {
	return vec.src
}

func (vec *Vector) SrcAt(i int) byte {
	return vec.src[i]
}

func (vec *Vector) SrcAddr() uint64 {
	return vec.addr
}

func (vec *Vector) Root() *Node {
	return vec.Get()
}

func (vec *Vector) Exists(key string) bool {
	n := vec.Root()
	return n.Exists(key)
}

func (vec *Vector) GetNode(depth int) (node *Node) {
	if vec.nodeL < len(vec.nodes) {
		node = &vec.nodes[vec.nodeL]
		node.Reset()
		vec.nodeL++
	} else {
		node = &Node{typ: TypeUnk}
		vec.nodes = append(vec.nodes, *node)
		vec.nodeL++
	}
	node.vecPtr, node.depth = vec.ptr(), depth
	node.key.vecPtr, node.val.vecPtr = node.vecPtr, node.vecPtr
	return
}

func (vec *Vector) GetNodeWT(depth int, typ Type) *Node {
	node := vec.GetNode(depth)
	node.typ = typ
	return node
}

func (vec *Vector) PutNode(idx int, node *Node) {
	vec.nodes[idx] = *node
}

func (vec *Vector) SetErrOffset(offset int) {
	vec.errOff = offset
}

func (vec *Vector) ErrorOffset() int {
	return vec.errOff
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
	vec.Index.reset()
}

func (vec *Vector) ptr() uintptr {
	if vec.selfPtr == 0 {
		vec.selfPtr = uintptr(unsafe.Pointer(vec))
	}
	return vec.selfPtr
}
