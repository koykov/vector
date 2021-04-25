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
		if err != nil || i >= n.limit {
			return nullNode
		}
		idx := vec.Index.val(n.depth+1, n.offset+i)
		child := &vec.nodes[idx]
		if len(keys[1:]) == 0 {
			return child
		} else {
			return n.Get(keys[1:]...)
		}
	}
	return nullNode
}

func (n *Node) GetObject(keys ...string) *Node {
	node := n.Get(keys...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

func (n *Node) GetArray(keys ...string) *Node {
	node := n.Get(keys...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

func (n *Node) GetBytes(keys ...string) []byte {
	node := n.Get(keys...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

func (n *Node) GetString(keys ...string) string {
	node := n.Get(keys...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

func (n *Node) GetBool(keys ...string) bool {
	node := n.Get(keys...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

func (n *Node) GetFloat(keys ...string) (float64, error) {
	node := n.Get(keys...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Float()
}

func (n *Node) GetInt(keys ...string) (int64, error) {
	node := n.Get(keys...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Int()
}

func (n *Node) GetUint(keys ...string) (uint64, error) {
	node := n.Get(keys...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Uint()
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
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

func (n *Node) GetArrayPS(path, sep string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

func (n *Node) GetBytesPS(path, sep string) []byte {
	vec := n.indirectVector()
	if vec == nil {
		return nil
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

func (n *Node) GetStringPS(path, sep string) string {
	vec := n.indirectVector()
	if vec == nil {
		return ""
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

func (n *Node) GetBoolPS(path, sep string) bool {
	vec := n.indirectVector()
	if vec == nil {
		return false
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

func (n *Node) GetFloatPS(path, sep string) (float64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Float()
}

func (n *Node) GetIntPS(path, sep string) (int64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Int()
}

func (n *Node) GetUintPS(path, sep string) (uint64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, sep, -1)
	node := n.Get(vec.bufSS...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Uint()
}
