package vector

import (
	"strconv"

	"github.com/koykov/bytealg"
)

// Get child node by given keys.
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
			if child.keyEqual(keys[0]) {
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

// Look and get child object by given keys.
func (n *Node) GetObject(keys ...string) *Node {
	node := n.Get(keys...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

// Look and get child array by given keys.
func (n *Node) GetArray(keys ...string) *Node {
	node := n.Get(keys...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

// Look and get child bytes by given keys.
func (n *Node) GetBytes(keys ...string) []byte {
	node := n.Get(keys...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

// Look and get child string by given keys.
func (n *Node) GetString(keys ...string) string {
	node := n.Get(keys...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

// Look and get child bool by given keys.
func (n *Node) GetBool(keys ...string) bool {
	node := n.Get(keys...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

// Look and get child float by given keys.
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

// Look and get child integer by given keys.
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

// Look and get child unsigned integer by given keys.
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

// Get child node by path and separator.
func (n *Node) GetPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS, path, separator, -1)
	return n.Get(vec.bufSS...)
}

// Look and get child object by given path and separator.
func (n *Node) GetObjectPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

// Look and get child array by given path and separator.
func (n *Node) GetArrayPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

// Look and get child bytes by given path and separator.
func (n *Node) GetBytesPS(path, separator string) []byte {
	vec := n.indirectVector()
	if vec == nil {
		return nil
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

// Look and get child string by given path and separator.
func (n *Node) GetStringPS(path, separator string) string {
	vec := n.indirectVector()
	if vec == nil {
		return ""
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

// Look and get child bool by given path and separator.
func (n *Node) GetBoolPS(path, separator string) bool {
	vec := n.indirectVector()
	if vec == nil {
		return false
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

// Look and get child float by given path and separator.
func (n *Node) GetFloatPS(path, separator string) (float64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Float()
}

// Look and get child integer by given path and separator.
func (n *Node) GetIntPS(path, separator string) (int64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Int()
}

// Look and get child unsigned int by given path and separator.
func (n *Node) GetUintPS(path, separator string) (uint64, error) {
	vec := n.indirectVector()
	if vec == nil {
		return 0, ErrInternal
	}
	vec.splitPath(path, separator)
	node := n.Get(vec.bufSS...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Uint()
}
