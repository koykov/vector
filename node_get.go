package vector

import (
	"strconv"

	"github.com/koykov/bytealg"
)

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
		for i := n.offset; i < n.limit; i++ {
			idx := vec.Index.val(n.depth+1, i)
			child := vec.nodes[idx]
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
		if err != nil || i >= n.limit {
			return nullNode
		}
		idx := vec.Index.val(n.depth+1, n.offset+i)
		child := vec.nodes[idx]
		if len(keys[1:]) == 0 {
			return child
		} else {
			return n.Get(keys[1:]...)
		}
	}
	return nullNode
}

func (n *Node) GetObject(keys ...string) *Node {
	c := n.Get(keys...)
	if c.Type() != TypeObj {
		return nullNode
	}
	return c.Object()
}

func (n *Node) GetArray(keys ...string) *Node {
	c := n.Get(keys...)
	if c.Type() != TypeArr {
		return nullNode
	}
	return c.Array()
}

func (n *Node) GetBytes(keys ...string) []byte {
	c := n.Get(keys...)
	if c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

func (n *Node) GetString(keys ...string) string {
	c := n.Get(keys...)
	if c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

func (n *Node) GetBool(keys ...string) bool {
	c := n.Get(keys...)
	if c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

func (n *Node) GetFloat(keys ...string) (float64, error) {
	c := n.Get(keys...)
	if c.typ == TypeNull {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Float()
}

func (n *Node) GetInt(keys ...string) (int64, error) {
	c := n.Get(keys...)
	if c.typ == TypeNull {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Int()
}

func (n *Node) GetUint(keys ...string) (uint64, error) {
	c := n.Get(keys...)
	if c.typ == TypeNull {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Uint()
}

func (n *Node) GetPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS, path, separator, -1)
	return n.Get(vec.bufSS...)
}

func (n *Node) GetObjectPS(path, sep string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.Type() != TypeObj {
		return nullNode
	}
	return c.Object()
}

func (n *Node) GetArrayPS(path, sep string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.Type() != TypeArr {
		return nullNode
	}
	return c.Array()
}

func (n *Node) GetBytesPS(path, sep string) []byte {
	vec := n.indirectVector()
	if vec == nil {
		return nil
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.Type() != TypeStr {
		return nil
	}
	return c.Bytes()
}

func (n *Node) GetStringPS(path, sep string) string {
	vec := n.indirectVector()
	if vec == nil {
		return ""
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.Type() != TypeStr {
		return ""
	}
	return c.String()
}

func (n *Node) GetBoolPS(path, sep string) bool {
	vec := n.indirectVector()
	if vec == nil {
		return false
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.Type() != TypeBool {
		return false
	}
	return c.Bool()
}

func (n *Node) GetFloatPS(path, sep string) (float64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.typ == TypeNull {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Float()
}

func (n *Node) GetIntPS(path, sep string) (int64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.typ == TypeNull {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Int()
}

func (n *Node) GetUintPS(path, sep string) (uint64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	c := n.Get(vec.bufSS...)
	if c.typ == TypeNull {
		return 0, ErrNotFound
	}
	if c.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return c.Uint()
}
