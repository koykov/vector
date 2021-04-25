package vector

import (
	"strconv"

	"github.com/koykov/bytealg"
)

func (vec *Vector) Get(keys ...string) *Node {
	if len(keys) == 0 {
		if vec.Len() > 0 {
			return &vec.nodes[0]
		}
		return nullNode
	}

	node := &vec.nodes[0]
	if node.typ != TypeObj && node.typ != TypeArr {
		if len(keys) > 1 {
			return nullNode
		}
		return node
	}

	if node.typ == TypeArr {
		return vec.getArr(node, keys...)
	}
	if node.typ == TypeObj {
		return vec.getObj(node, keys...)
	}
	return node
}

func (vec *Vector) GetObject(keys ...string) *Node {
	node := vec.Get(keys...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

func (vec *Vector) GetArray(keys ...string) *Node {
	node := vec.Get(keys...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

func (vec *Vector) GetBytes(keys ...string) []byte {
	node := vec.Get(keys...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

func (vec *Vector) GetString(keys ...string) string {
	node := vec.Get(keys...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

func (vec *Vector) GetBool(keys ...string) bool {
	node := vec.Get(keys...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

func (vec *Vector) GetFloat(keys ...string) (float64, error) {
	node := vec.Get(keys...)
	if node.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Float()
}

func (vec *Vector) GetInt(keys ...string) (int64, error) {
	node := vec.Get(keys...)
	if node.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Int()
}

func (vec *Vector) GetUint(keys ...string) (uint64, error) {
	node := vec.Get(keys...)
	if node.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Uint()
}

func (vec *Vector) GetPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	return vec.Get(vec.bufSS...)
}

func (vec *Vector) GetObjectPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

func (vec *Vector) GetArrayPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

func (vec *Vector) GetBytesPS(path, separator string) []byte {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

func (vec *Vector) GetStringPS(path, separator string) string {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

func (vec *Vector) GetBoolPS(path, separator string) bool {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

func (vec *Vector) GetFloatPS(path, separator string) (float64, error) {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Float()
}

func (vec *Vector) GetIntPS(path, separator string) (int64, error) {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Int()
}

func (vec *Vector) GetUintPS(path, separator string) (uint64, error) {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	node := vec.Get(vec.bufSS...)
	if node.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Uint()
}

func (vec *Vector) getArr(root *Node, keys ...string) *Node {
	if len(keys) == 0 {
		return root
	}
	k, err := strconv.Atoi(keys[0])
	if err != nil || k >= root.limit {
		return nullNode
	}
	i := vec.Index.val(root.depth+1, root.offset+k)
	node := &vec.nodes[i]
	tail := keys[1:]
	if node.typ != TypeArr && node.typ != TypeObj {
		if len(tail) > 0 {
			return nullNode
		}
		return node
	}
	if node.typ == TypeArr {
		return vec.getArr(node, tail...)
	}
	if node.typ == TypeObj {
		return vec.getObj(node, tail...)
	}
	return nullNode
}

func (vec *Vector) getObj(root *Node, keys ...string) *Node {
	if len(keys) == 0 {
		return root
	}
	var node *Node
	for i := root.offset; i < root.limit; i++ {
		k := vec.Index.val(root.depth+1, i)
		node = &vec.nodes[k]
		if node.key.String() == keys[0] {
			break
		}
	}
	if node == nil {
		return nullNode
	}
	tail := keys[1:]
	if node.typ != TypeArr && node.typ != TypeObj {
		if len(tail) > 0 {
			return nullNode
		}
		return node
	}
	if node.typ == TypeArr {
		return vec.getArr(node, tail...)
	}
	if node.typ == TypeObj {
		return vec.getObj(node, tail...)
	}
	return nullNode
}
