package vector

import (
	"io"
	"reflect"
	"strings"
	"unsafe"

	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
)

// Vector parser object.
type Vector struct {
	bitset.Bitset
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

// Parse parses source bytes.
func (vec *Vector) Parse(_ []byte) error {
	return ErrNotImplement
}

// ParseStr parses source string.
func (vec *Vector) ParseStr(_ string) error {
	return ErrNotImplement
}

// ParseCopy copies source bytes and parse it.
func (vec *Vector) ParseCopy(_ []byte) error {
	return ErrNotImplement
}

// ParseCopyStr copies source string and parse it.
func (vec *Vector) ParseCopyStr(_ string) error {
	return ErrNotImplement
}

// Beautify formats first root node in human-readable representation.
//
// Second and next roots must beautify manually by call Beautify method of each node.
func (vec *Vector) Beautify(w io.Writer) error {
	if vec.Helper == nil {
		return ErrNoHelper
	}
	root := vec.Root()
	return vec.Helper.Beautify(w, root)
}

// Marshal serializes first root node.
//
// Second and next roots must beautify manually by call Marshal method of each node.
func (vec *Vector) Marshal(w io.Writer) error {
	if vec.Helper == nil {
		return ErrNoHelper
	}
	root := vec.Root()
	return vec.Helper.Marshal(w, root)
}

// SetSrc sets source bytes.
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

// Len returns length of nodes array.
func (vec *Vector) Len() int {
	return vec.nodeL
}

// SrcLen returns length of source bytes.
func (vec *Vector) SrcLen() int {
	return len(vec.src)
}

// Src returns raw source bytes.
//
// Please note, source bytes may be "corrupt" when unescaped.
func (vec *Vector) Src() []byte {
	return vec.src
}

// SrcAt returns byte at position i.
func (vec *Vector) SrcAt(i int) byte {
	return vec.src[i]
}

// SrcAddr returns source address in virtual memory.
func (vec *Vector) SrcAddr() uint64 {
	return vec.addr
}

// Root returns root node.
func (vec *Vector) Root() *Node {
	return vec.Get()
}

// RootLen returns count of root nodes.
func (vec *Vector) RootLen() int {
	return vec.Index.Len(0)
}

// RootByIndex returns root node by given index.
//
// For cases when one vector instance uses for parse many sources.
func (vec *Vector) RootByIndex(index int) *Node {
	rootRow := vec.Index.GetRow(0)
	if index < 0 || index >= len(rootRow) {
		return nullNode
	}
	idx := rootRow[index]
	return &vec.nodes[idx]
}

// RootTop returns last root node.
func (vec *Vector) RootTop() *Node {
	rootRow := vec.Index.GetRow(0)
	if len(rootRow) == 0 {
		return nullNode
	}
	return &vec.nodes[rootRow[len(rootRow)-1]]
}

// Each applies custom function to each root node.
func (vec *Vector) Each(fn func(idx int, node *Node)) {
	rootRow := vec.Index.GetRow(0)
	l := len(rootRow)
	if l == 0 {
		return
	}
	c := 0
	_ = rootRow[l-1]
	for i := 0; i < l; i++ {
		root := &vec.nodes[i]
		fn(c, root)
		c++
	}
}

// RemoveIf deletes all root nodes satisfies condition cond.
func (vec *Vector) RemoveIf(cond func(idx int, node *Node) bool) {
	idx := vec.Index.GetRow(0)
	l := len(idx)
	if l == 0 {
		return
	}
	c := 0
	_ = idx[l-1]
	for i := 0; i < len(idx); {
		i_ := idx[i]
		cn := &vec.nodes[i_]
		if cond(c, cn) {
			idx = append(idx[:i], idx[i+1:]...)
			c++
			continue
		}
		i++
		c++
	}
}

// Exists checks if node exists by given key.
func (vec *Vector) Exists(key string) bool {
	n := vec.Root()
	return n.Exists(key)
}

// GetNode returns new node with index from position matrix by given depth.
func (vec *Vector) GetNode(depth int) (node *Node, idx int) {
	node, idx = vec.getNode(depth)
	vec.Index.Register(depth, idx)
	return
}

// GetNodeWT returns node and set type at once.
func (vec *Vector) GetNodeWT(depth int, typ Type) (*Node, int) {
	node, idx := vec.GetNode(depth)
	node.typ = typ
	return node, idx
}

// NodeAt returns node at given position.
func (vec *Vector) NodeAt(idx int) *Node {
	if idx < 0 || idx >= vec.Len() {
		return nullNode
	}
	return &vec.nodes[idx]
}

// GetChild get node and register it as a child of root node.
//
// Similar to GetNode.
func (vec *Vector) GetChild(root *Node, depth int) (*Node, int) {
	return vec.GetChildWT(root, depth, TypeUnk)
}

// GetChildWT get node, register it as a child of root node and set type at once.
func (vec *Vector) GetChildWT(root *Node, depth int, typ Type) (*Node, int) {
	node, idx := vec.getNode(depth)
	node.typ = typ
	node.pptr = root.ptr()
	root.SetLimit(vec.Index.Register(depth, idx))
	return node, idx
}

// Generic node getter.
//
// Returns new node and its index in the nodes array.
func (vec *Vector) getNode(depth int) (node *Node, idx int) {
	n := len(vec.nodes)
	if vec.nodeL < n {
		_ = vec.nodes[n-1]
		node = &vec.nodes[vec.nodeL]
		node.Reset()
	} else {
		vec.nodes = append(vec.nodes, Node{typ: TypeUnk})
		n++
		_ = vec.nodes[n-1]
		node = &vec.nodes[n-1]
	}
	vec.nodeL++
	node.vptr, node.depth = vec.ptr(), depth
	node.key.vptr, node.val.vptr = node.vptr, node.vptr
	idx = vec.Len() - 1
	node.idx = idx
	return
}

// PutNode returns node back to the array.
func (vec *Vector) PutNode(idx int, node *Node) {
	l := unsafe.Pointer(&vec.nodes[idx])
	r := unsafe.Pointer(node)
	if uintptr(l) != uintptr(r) {
		vec.nodes[idx] = *node
	}
}

// SetErrOffset sets error offset.
func (vec *Vector) SetErrOffset(offset int) {
	vec.errOff = offset
}

// ErrorOffset returns last error offset.
func (vec *Vector) ErrorOffset() int {
	return vec.errOff
}

// Reset vector data.
func (vec *Vector) Reset() {
	if vec.nodeL == 0 {
		return
	}
	_ = vec.nodes[vec.nodeL-1]
	vec.buf, vec.src = vec.buf[:0], nil
	vec.bufSS = vec.bufSS[:0]
	vec.addr, vec.nodeL, vec.errOff = 0, 0, 0
	vec.Index.reset()
	vec.Bitset.Reset()
}

// ForgetFrom forgets nodes from given position to the end of the array.
func (vec *Vector) ForgetFrom(idx int) {
	if idx >= vec.nodeL {
		return
	}
	for i := idx; i < vec.nodeL; i++ {
		vec.nodes[i].Reset()
	}
	vec.nodeL = idx
}

// KeepPtr guarantees that vector object wouldn't be collected by GC.
//
// Typically, vector objects uses together with pools and GC doesn't collect them. But for cases like
// vec := &Vector{...}
// node := vec.Get("foo", "bar")
// <- here GC may collect vec
// fmt.Println(node.String()) <- invalid operation due to vec already has been collected
// vec.KeepPtr() <- just call me to avoid that trouble
func (vec *Vector) KeepPtr() {
	_ = vec.ptr()
}

// Return self pointer of the vector.
func (vec *Vector) ptr() uintptr {
	if vec.selfPtr == 0 {
		vec.selfPtr = uintptr(unsafe.Pointer(vec))
	}
	return vec.selfPtr
}

// Split path by given separator.
//
// Caution! Don't user "@" as a separator, it will break work with attributes.
// TODO: consider escaped at symbol "\@".
func (vec *Vector) splitPath(path, separator string) {
	vec.bufSS = bytealg.AppendSplit(vec.bufSS[:0], path, separator, -1)
	ti := len(vec.bufSS) - 1
	if ti < 0 {
		return
	}
	tail := vec.bufSS[ti]
	if p := strings.IndexByte(tail, '@'); p != -1 {
		if p > 0 {
			if len(tail[p:]) > 1 {
				vec.bufSS = append(vec.bufSS, tail[p:])
			}
			vec.bufSS[ti] = vec.bufSS[ti][:p]
		}
	}
}
