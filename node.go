package vector

import (
	"strconv"
	"unsafe"

	"github.com/koykov/bytealg"
)

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

type Node struct {
	typ      Type
	key, val Byteptr

	depth, offset, limit int

	vecPtr uintptr
}

var (
	nullNode = &Node{typ: TypeNull}
)

func (n *Node) SetType(typ Type) {
	n.typ = typ
}

func (n *Node) Type() Type {
	return n.typ
}

func (n *Node) Depth() int {
	return n.depth
}

func (n *Node) SetOffset(offset int) {
	n.offset = offset
}

func (n *Node) Offset() int {
	return n.offset
}

func (n *Node) SetLimit(limit int) {
	n.limit = limit
}

func (n *Node) Limit() int {
	if n.limit != n.offset && n.limit >= n.offset {
		return n.limit - n.offset
	} else if n.limit == 0 && n.offset > 0 {
		return 0
	}
	return 1
}

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
		c := vec.nodes[k]
		if c.key.String() == key {
			return true
		}
	}
	return false
}

func (n *Node) Object() *Node {
	if n.typ != TypeObj {
		return nullNode
	}
	return n
}

func (n *Node) Array() *Node {
	if n.typ != TypeArr {
		return nullNode
	}
	return n
}

func (n *Node) Key() *Byteptr {
	return &n.key
}

func (n *Node) KeyBytes() []byte {
	if n.key.Offset() != 0 && n.key.Limit() > 0 {
		return n.key.RawBytes()
	}
	return nil
}

func (n *Node) KeyString() string {
	if n.key.Offset() != 0 && n.key.Limit() > 0 {
		return n.key.String()
	}
	return ""
}

func (n *Node) Value() *Byteptr {
	return &n.val
}

func (n *Node) Bytes() []byte {
	if n.typ != TypeStr && n.typ != TypeNum && n.typ != TypeBool && n.typ != TypeAttr {
		return nil
	}
	return n.val.Bytes()
}

func (n *Node) ForceBytes() []byte {
	return n.val.Bytes()
}

func (n *Node) RawBytes() []byte {
	return n.val.Bytes()
}

func (n *Node) String() string {
	if n.typ != TypeStr && n.typ != TypeNum && n.typ != TypeBool && n.typ != TypeAttr {
		return ""
	}
	return n.val.String()
}

func (n *Node) ForceString() string {
	return n.val.String()
}

func (n *Node) Bool() bool {
	if n.typ != TypeBool {
		return false
	}
	return bytealg.ToLowerStr(n.val.String()) == "true"
}

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

func (n *Node) Each(fn func(idx int, node *Node)) {
	idx := n.childs()
	vec := n.indirectVector()
	if len(idx) == 0 || vec == nil {
		return
	}
	c := 0
	for _, i := range idx {
		cn := &vec.nodes[i]
		fn(c, cn)
		c++
	}
}

func (n *Node) Look(key string) *Node {
	if n.typ != TypeObj {
		return nullNode
	}
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	ci := n.childs()
	for _, i := range ci {
		c := &vec.nodes[i]
		if key == c.key.String() {
			return c
		}
	}
	return nullNode
}

func (n *Node) At(idx int) *Node {
	if n.typ != TypeArr {
		return nullNode
	}
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	ci := n.childs()
	h := -1
	for _, i := range ci {
		if i == idx {
			h = i
			break
		}
	}
	if h >= 0 {
		return &vec.nodes[h]
	}
	return nil
}

func (n *Node) Reset() {
	n.typ = TypeUnk
	n.key.Reset()
	n.val.Reset()
	n.depth, n.offset, n.limit, n.vecPtr = 0, 0, 0, 0
}

func (n *Node) childs() []int {
	if vec := n.indirectVector(); vec != nil {
		var e = n.limit
		if e == 0 {
			e = n.offset + 1
		}
		return vec.Index.get(n.depth+1, n.offset, e)
	}
	return nil
}

func (n *Node) indirectVector() *Vector {
	if n.vecPtr == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(n.vecPtr))
}
