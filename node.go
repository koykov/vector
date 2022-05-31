package vector

import (
	"strconv"
	"unsafe"

	"github.com/koykov/bytealg"
)

// Node type.
type Type int

const (
	TypeUnk Type = iota
	TypeNull
	TypeObj
	TypeArr
	TypeStr
	TypeNum
	TypeBool
	TypeAttr
)

// Node object.
type Node struct {
	// Node type.
	typ Type
	// Key/value byteptr objects.
	key, val, aka Byteptr
	// Node index in array, depth in a index tree, offset in index row and limit of childs in index row.
	idx, depth, offset, limit int
	// Raw pointers to vector and parent node.
	// It's safe to usi uintptr here because vector guaranteed to exist while the node is alive and isn't garbage
	// collected.
	vptr, pptr uintptr
}

var (
	// Null node instance. Will return for empty results.
	nullNode = &Node{typ: TypeNull}
)

// Set type of the node.
func (n *Node) SetType(typ Type) {
	n.typ = typ
}

// Get node type.
func (n *Node) Type() Type {
	return n.typ
}

// Get node index in the array of nodes.
func (n *Node) Index() int {
	return n.idx
}

// Get node depth in index tree.
func (n *Node) Depth() int {
	return n.depth
}

// Set offset in the index row.
func (n *Node) SetOffset(offset int) *Node {
	n.offset = offset
	return n
}

// Get offset of childs in the index row.
func (n *Node) Offset() int {
	return n.offset
}

// Set limit of childs in index row.
func (n *Node) SetLimit(limit int) *Node {
	n.limit = limit
	return n
}

// Get limit of childs in index row.
func (n *Node) Limit() int {
	if n.limit != n.offset && n.limit >= n.offset {
		return n.limit - n.offset
	} else if n.limit == 0 && n.offset > 0 {
		return 0
	}
	return 1
}

// Check if given key exists in the node.
func (n *Node) Exists(key string) bool {
	if n.typ != TypeObj {
		return false
	}
	vec := n.indirectVector()
	if vec == nil {
		return false
	}
	for i := n.offset; i < n.limit; i++ {
		k := vec.Index.val(n.depth+1, i)
		c := &vec.nodes[k]
		if c.key.String() == key {
			return true
		}
	}
	return false
}

// Check node is object and return it.
func (n *Node) Object() *Node {
	if n.typ != TypeObj {
		return nullNode
	}
	return n
}

// Check node is array and return it.
func (n *Node) Array() *Node {
	if n.typ != TypeArr {
		return nullNode
	}
	return n
}

// Get key as byteptr object.
func (n *Node) Key() *Byteptr {
	return &n.key
}

// Get key as bytes.
func (n *Node) KeyBytes() []byte {
	return n.key.RawBytes()
}

// Get key as string.
func (n *Node) KeyString() string {
	return n.key.String()
}

// Get value as byteptr object.
func (n *Node) Value() *Byteptr {
	return &n.val
}

// Get value as bytes.
//
// Allow only for [string, number, bool, attribute] types.
func (n *Node) Bytes() []byte {
	if n.typ != TypeStr && n.typ != TypeNum && n.typ != TypeBool && n.typ != TypeAttr {
		return nil
	}
	return n.val.Bytes()
}

// Get value as bytes independent of the type.
func (n *Node) ForceBytes() []byte {
	return n.val.Bytes()
}

// Get value as bytes without implement any conversion logic.
func (n *Node) RawBytes() []byte {
	return n.val.RawBytes()
}

// Get value as string.
//
// Allow only for [string, number, bool, attribute] types.
func (n *Node) String() string {
	if n.typ != TypeStr && n.typ != TypeNum && n.typ != TypeBool && n.typ != TypeAttr {
		return ""
	}
	return n.val.String()
}

// Get value as string independent of the type.
func (n *Node) ForceString() string {
	return n.val.String()
}

// AKA returns pointer to alias object.
func (n *Node) AKA() *Byteptr {
	return &n.aka
}

// Get value as boolean.
func (n *Node) Bool() bool {
	if n.typ != TypeBool {
		return false
	}
	return bytealg.ToLowerStr(n.val.String()) == "true"
}

// Get value as float number.
func (n *Node) Float() (float64, error) {
	if n.typ != TypeNum {
		return 0, ErrIncompatType
	}
	f, err := strconv.ParseFloat(n.val.String(), 64)
	if err != nil {
		return 0, err
	}
	return f, nil
}

// Get value as integer.
func (n *Node) Int() (int64, error) {
	if n.typ != TypeNum {
		return 0, ErrIncompatType
	}
	i, err := strconv.ParseInt(n.val.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return i, nil
}

// Get value as unsigned integer.
func (n *Node) Uint() (uint64, error) {
	if n.typ != TypeNum {
		return 0, ErrIncompatType
	}
	u, err := strconv.ParseUint(n.val.String(), 10, 64)
	if err != nil {
		return 0, err
	}
	return u, nil
}

// Apply custom function to each child of the node.
func (n *Node) Each(fn func(idx int, node *Node)) {
	idx := n.childrenIdx()
	vec := n.indirectVector()
	if len(idx) == 0 || vec == nil {
		return
	}
	c := 0
	_ = idx[len(idx)-1]
	for _, i := range idx {
		cn := &vec.nodes[i]
		fn(c, cn)
		c++
	}
}

// Look for the child node by given key.
//
// May be used only for object nodes.
func (n *Node) Look(key string) *Node {
	if n.typ != TypeObj {
		return nullNode
	}
	ci := n.childrenIdx()
	vec := n.indirectVector()
	if len(ci) == 0 || vec == nil {
		return nullNode
	}
	_ = ci[len(ci)-1]
	for _, i := range ci {
		c := &vec.nodes[i]
		if key == c.key.String() {
			return c
		}
	}
	return nullNode
}

// Get node from array at position idx.
//
// May be used only for array nodes.
func (n *Node) At(idx int) *Node {
	if n.typ != TypeArr {
		return nullNode
	}
	ci := n.childrenIdx()
	vec := n.indirectVector()
	if len(ci) == 0 || len(ci) <= idx || ci[idx] == 0 || vec == nil {
		return nullNode
	}
	return &vec.nodes[ci[idx]]
}

// Reset the node.
func (n *Node) Reset() *Node {
	n.typ = TypeUnk
	n.key.Reset()
	n.val.Reset()
	n.aka.Reset()
	n.depth, n.offset, n.limit, n.vptr = 0, 0, 0, 0
	return n
}

// Get list of children indexes.
func (n *Node) childrenIdx() []int {
	if vec := n.indirectVector(); vec != nil {
		var limit = n.limit
		if limit == 0 {
			limit = n.offset + 1
		}
		return vec.Index.get(n.depth+1, n.offset, limit)
	}
	return nil
}

// Get list of children nodes.
func (n *Node) Children() []Node {
	if ci := n.childrenIdx(); len(ci) > 0 {
		if vec := n.indirectVector(); vec != nil {
			offset, limit := ci[0], ci[len(ci)-1]+1
			if limit >= offset && limit <= vec.nodeL {
				return vec.nodes[offset:limit]
			}
		}
	}
	return nil
}

// Swap node with another given node in the nodes array.
func (n *Node) SwapWith(node *Node) {
	if vec := n.indirectVector(); vec != nil {
		i, j := n.idx, node.idx
		if i < vec.nodeL && j < vec.nodeL {
			vec.nodes[i].idx, vec.nodes[j].idx = j, i
			vec.nodes[i], vec.nodes[j] = vec.nodes[j], vec.nodes[i]
		}
	}
}

// Check key equality.
//
// Also check node type for keys with "@" prefix (must be TypeAttr).
func (n *Node) keyEqual(key string) bool {
	if key[0] == '@' {
		key = key[1:]
		return n.key.String() == key && n.typ == TypeAttr
	}
	return n.typ != TypeAttr && n.key.String() == key
}

// Return self pointer of the node.
func (n *Node) ptr() uintptr {
	return uintptr(unsafe.Pointer(n))
}

// Restore the entire vector object from the unsafe pointer.
//
// This needs to reduce pointers count and avoids redundant GC checks.
func (n *Node) indirectVector() *Vector {
	if n.vptr == 0 {
		return nil
	}
	return indirectVector1(n.vptr)
}

// Restore the entire node object from the unsafe pointer.
//
// This needs to reduce pointers count and avoids redundant GC checks.
func (n *Node) indirectNode() *Node {
	if n.pptr == 0 {
		return nil
	}
	return indirectNode1(n.pptr)
}
