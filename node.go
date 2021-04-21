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
	key, val bytealg.Byteptr

	depth, offset, length int

	vecPtr uintptr
}

var (
	nullNode = &Node{typ: TypeNull}
)

func (n *Node) Type() Type {
	return n.typ
}

func (n *Node) Len() int {
	if n.length != n.offset && n.length >= n.offset {
		return n.length - n.offset
	}
	return 1
}

func (n *Node) Get(keys ...string) *Node {
	if len(keys) == 0 {
		return n
	}
	if n.typ != TypeObj && n.typ != TypeArr {
		return nullNode
	}
	var vec *Vector
	if vec = n.indirectVector(); vec == nil {
		return n
	}
	if n.typ == TypeObj {
		for i := n.offset; i < n.length; i++ {
			idx := vec.index.val(n.depth+1, i)
			child := &vec.nodes[idx]
			if child.key.String() == keys[0] {
				if len(keys[1:]) == 0 {
					return child
				} else {
					return child.Get(keys[1:]...)
				}
			}
		}
	}
	if n.typ == TypeArr {
		i, err := strconv.Atoi(keys[0])
		if err != nil || i >= n.Len() {
			return nullNode
		}
		idx := vec.index.val(n.depth+1, n.offset+i)
		child := &vec.nodes[idx]
		if len(keys[1:]) == 0 {
			return child
		} else {
			return n.Get(keys[1:]...)
		}
	}
	return nullNode
}

func (n *Node) GetPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS, path, separator, -1)
	return n.Get(vec.bufSS...)
}

func (n *Node) Dot(path string) *Node {
	return n.GetPS(path, ".")
}

func (n *Node) Exists(key string) bool {
	if n.typ != TypeObj {
		return false
	}
	vec := n.indirectVector()
	if vec == nil {
		return false
	}
	for i := n.offset; i < n.length; i++ {
		k := vec.index.val(n.depth+1, i)
		c := &vec.nodes[k]
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

func (n *Node) Key() []byte {
	if n.key.Offset() != 0 && n.key.Len() > 0 {
		return n.key.Bytes()
	}
	return nil
}

func (n *Node) KeyString() string {
	if n.key.Offset() != 0 && n.key.Len() > 0 {
		return n.key.String()
	}
	return ""
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
		c := vec.nodes[i]
		if key == c.key.String() {
			return &c
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
	n.key.Set(0, 0)
	n.val.Set(0, 0)
	n.depth, n.offset, n.length, n.vecPtr = 0, 0, 0, 0
}

func (n *Node) childs() []int {
	if vec := n.indirectVector(); vec != nil {
		var e = n.length
		if e == 0 {
			e = n.offset + 1
		}
		return vec.index.get(n.depth+1, n.offset, e)
	}
	return nil
}

func (n *Node) indirectVector() *Vector {
	if n.vecPtr == 0 {
		return nil
	}
	return (*Vector)(unsafe.Pointer(n.vecPtr))
}
