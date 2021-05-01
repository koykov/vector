package vector

import (
	"io"
	"reflect"
	"unsafe"
)

// Vector parser object.
type Vector struct {
	// Source data to parse.
	src []byte
	// Source data pointer.
	addr uint64
	// Buffers.
	buf   []byte
	bufSS []string
	// Self pointer.
	selfPtr uintptr
	// List of nodes and length of it.
	nodes []Node
	nodeL int
	// Error offset.
	errOff int
	// Nodes index.
	Index Index
	// External helper object.
	Helper Helper
}

// Helper object interface.
type Helper interface {
	// Convert byteptr to byte slice and apply custom logic.
	ConvertByteptr(*Byteptr) []byte
}

// Parse source bytes.
func (vec *Vector) Parse(_ []byte) error {
	return ErrNotImplement
}

// Parse source string.
func (vec *Vector) ParseStr(_ string) error {
	return ErrNotImplement
}

// Copy source bytes and parse it.
func (vec *Vector) ParseCopy(_ []byte) error {
	return ErrNotImplement
}

// Copy source string and parse it.
func (vec *Vector) ParseCopyStr(_ string) error {
	return ErrNotImplement
}

// Format vector in human readable representation.
func (vec *Vector) Beautify(_ io.Writer) error {
	return ErrNotImplement
}

// Set source bytes to parse.
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
	// Get source data address.
	h := (*reflect.SliceHeader)(unsafe.Pointer(&vec.src))
	vec.addr = uint64(h.Data)
	return nil
}

// Get length of nodes array.
func (vec *Vector) Len() int {
	return vec.nodeL
}

// Get length of source bytes.
func (vec *Vector) SrcLen() int {
	return len(vec.src)
}

// Get raw source bytes.
func (vec *Vector) Src() []byte {
	return vec.src
}

// Get byte at position i.
func (vec *Vector) SrcAt(i int) byte {
	return vec.src[i]
}

// Get source address in virtual memory.
func (vec *Vector) SrcAddr() uint64 {
	return vec.addr
}

// Get root node.
func (vec *Vector) Root() *Node {
	return vec.Get()
}

// Check if node exists by given key.
func (vec *Vector) Exists(key string) bool {
	n := vec.Root()
	return n.Exists(key)
}

// Get node from the nodes array.
func (vec *Vector) GetNode(depth int) (node *Node, idx int) {
	node, idx = vec.getNode(depth)
	vec.Index.Register(depth, idx)
	return
}

// Get node and set type at once.
func (vec *Vector) GetNodeWT(depth int, typ Type) (*Node, int) {
	node, idx := vec.GetNode(depth)
	node.typ = typ
	return node, idx
}

// Get node and register it as a child of root node.
func (vec *Vector) GetChild(root *Node, depth int) (*Node, int) {
	node, idx := vec.getNode(depth)
	root.SetLimit(vec.Index.Register(depth, idx))
	return node, idx
}

// Get node,  register it as a child of root node and set type at once.
func (vec *Vector) GetChildWT(root *Node, depth int, typ Type) (*Node, int) {
	node, idx := vec.getNode(depth)
	node.typ = typ
	root.SetLimit(vec.Index.Register(depth, idx))
	return node, idx
}

// Generic node getter.
//
// Returns new node and its index in the nodes array.
func (vec *Vector) getNode(depth int) (node *Node, idx int) {
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
	idx = vec.Len() - 1
	return
}

// Return node back to the array.
func (vec *Vector) PutNode(idx int, node *Node) {
	vec.nodes[idx] = *node
}

// Set error offset.
func (vec *Vector) SetErrOffset(offset int) {
	vec.errOff = offset
}

// Get last error offset.
func (vec *Vector) ErrorOffset() int {
	return vec.errOff
}

// Reset vector data.
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

// Return self pointer of the vector.
func (vec *Vector) ptr() uintptr {
	if vec.selfPtr == 0 {
		vec.selfPtr = uintptr(unsafe.Pointer(vec))
	}
	return vec.selfPtr
}
