package vector

import (
	"io"
	"os"
	"unicode/utf8"
	"unsafe"

	"github.com/koykov/bitset"
	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/entry"
	"github.com/koykov/openrt"
)

// Vector parser object.
type Vector struct {
	bitset.Bitset
	// Source data to parse.
	src []byte
	// Source data pointer.
	addr uintptr
	// Buffers.
	buf   []byte
	bufKE []entry.Entry64
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

// ParseString parses source string.
func (vec *Vector) ParseString(_ string) error {
	return ErrNotImplement
}

// ParseCopy copies source bytes and parse it.
func (vec *Vector) ParseCopy(_ []byte) error {
	return ErrNotImplement
}

// ParseCopyString copies source string and parse it.
func (vec *Vector) ParseCopyString(_ string) error {
	return ErrNotImplement
}

// ParseFile reads file contents and parse it.
func (vec *Vector) ParseFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return vec.ParseReader(f)
}

// ParseReader reads source from r and parse it.
func (vec *Vector) ParseReader(r io.Reader) (err error) {
	const bufsz = 512
	for {
		off := len(vec.buf)
		vec.buf = bytealg.GrowDelta(vec.buf, bufsz)
		var n int
		n, err = r.Read(vec.buf[off:])
		if err == io.EOF || n < bufsz {
			vec.buf = vec.buf[:off+n]
			err = nil
			break
		}
		if err != nil {
			return err
		}
	}
	// Each submodule must provide own implementation. So base method always return "not implement" error.
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
	h := (*byteconv.SliceHeader)(unsafe.Pointer(&vec.src))
	vec.addr = h.Data
	vec.selfPtr = uintptr(unsafe.Pointer(vec))
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

// ReadRuneAt returns rune at position i.
// Please note, it's your responsibility to specify right position `i`.
func (vec *Vector) ReadRuneAt(i int) (r rune, w int, err error) {
	if i < 0 || i >= len(vec.src) {
		err = io.ErrUnexpectedEOF
		return
	}
	r, w = utf8.DecodeRune(vec.src[i:])
	return
}

// SrcAddr returns source address in virtual memory.
func (vec *Vector) SrcAddr() uintptr {
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

// AcquireNode allocates new node on given depth and returns it with index.
func (vec *Vector) AcquireNode(depth int) (node *Node, idx int) {
	return vec.AcquireNodeWithType(depth, TypeUnknown)
}

// AcquireNodeWithType allocates new node on given depth and returns it with index.
// Similar to AcquireNode but sets node's type at once.
func (vec *Vector) AcquireNodeWithType(depth int, typ Type) (*Node, int) {
	node, idx := vec.ackNode(depth)
	vec.Index.Register(depth, idx)
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

// AcquireChild allocates new node and marks it as a child of root node.
//
// Similar to AcquireNode.
func (vec *Vector) AcquireChild(root *Node, depth int) (*Node, int) {
	return vec.AcquireChildWithType(root, depth, TypeUnknown)
}

// AcquireChildWithType allocates new node, mark it as child of root and set type at once.
func (vec *Vector) AcquireChildWithType(root *Node, depth int, typ Type) (*Node, int) {
	node, idx := vec.ackNode(depth)
	node.typ = typ
	node.pptr = root.ptr()
	root.SetLimit(vec.Index.Register(depth, idx))
	return node, idx
}

// Generic node getter.
//
// Returns new node and its index in the nodes array.
func (vec *Vector) ackNode(depth int) (*Node, int) {
	n := len(vec.nodes)
	if n > 0 {
		_ = vec.nodes[n-1]
	}
	var node *Node
	if vec.nodeL < n {
		node = &vec.nodes[vec.nodeL]
	} else {
		vec.nodes = append(vec.nodes, Node{typ: TypeUnknown})
		node = &vec.nodes[n]
	}
	node.depth = depth
	node.vptr = vec.selfPtr
	node.key.vptr = node.vptr
	node.val.vptr = node.vptr
	node.idx = vec.nodeL
	vec.nodeL++
	return node, node.idx
}

// ReleaseNode returns node back to the vector.
func (vec *Vector) ReleaseNode(idx int, node *Node) {
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

// Prealloc prepares space for further parse.
func (vec *Vector) Prealloc(size uint) {
	if ul := uint(len(vec.nodes)); ul < size {
		vec.nodes = append(vec.nodes, make([]Node, size-ul+1)...)
	}
}

// Reset vector data.
func (vec *Vector) Reset() {
	if vec.nodeL == 0 {
		return
	}
	if !vec.CheckBit(FlagNoClear) {
		openrt.MemclrUnsafe(unsafe.Pointer(&vec.nodes[0]), vec.nodeL*nodeSize)
	}
	vec.nodeL = 0

	vec.buf, vec.src = vec.buf[:0], nil
	vec.bufKE = vec.bufKE[:0]
	vec.addr, vec.nodeL, vec.errOff = 0, 0, 0
	vec.Index.reset()
	vec.Bitset.Reset()
	vec.Bitset.SetBit(FlagInit, vec.Helper != nil)
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
	return uintptr(unsafe.Pointer(vec))
}
