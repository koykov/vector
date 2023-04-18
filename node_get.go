package vector

import (
	"strconv"

	"github.com/koykov/bytealg"
)

// Get returns child node by given keys.
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
				tail := keys[1:]
				if len(tail) == 0 {
					return child
				} else {
					return child.Get(tail...)
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
		tail := keys[1:]
		if len(tail) == 0 {
			return child
		} else if n.aka.Len() > 0 && len(tail) == 1 && n.aka.String() == tail[0] {
			return n
		} else {
			return child.Get(tail...)
		}
	}
	return nullNode
}

// GetObject looks and get child object by given keys.
func (n *Node) GetObject(keys ...string) *Node {
	node := n.Get(keys...)
	if node.Type() != TypeObj {
		return nullNode
	}
	return node.Object()
}

// GetArray looks and get child array by given keys.
func (n *Node) GetArray(keys ...string) *Node {
	node := n.Get(keys...)
	if node.Type() != TypeArr {
		return nullNode
	}
	return node.Array()
}

// GetBytes looks and get child bytes by given keys.
func (n *Node) GetBytes(keys ...string) []byte {
	node := n.Get(keys...)
	if node.Type() != TypeStr {
		return nil
	}
	return node.Bytes()
}

// GetString looks and get child string by given keys.
func (n *Node) GetString(keys ...string) string {
	node := n.Get(keys...)
	if node.Type() != TypeStr {
		return ""
	}
	return node.String()
}

// GetBool looks and get child bool by given keys.
func (n *Node) GetBool(keys ...string) bool {
	node := n.Get(keys...)
	if node.Type() != TypeBool {
		return false
	}
	return node.Bool()
}

// GetFloat looks and get child float by given keys.
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

// GetInt looks and get child integer by given keys.
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

// GetUint looks and get child unsigned integer by given keys.
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

// GetPS returns child node by path and separator.
func (n *Node) GetPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	vec.bufSS = bytealg.AppendSplit(vec.bufSS[:0], path, separator, -1)
	return n.Get(vec.bufSS...)
}

// GetObjectPS looks and get child object by given path and separator.
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

// GetArrayPS looks and get child array by given path and separator.
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

// GetBytesPS looks and get child bytes by given path and separator.
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

// GetStringPS looks and get child string by given path and separator.
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

// GetBoolPS looks and get child bool by given path and separator.
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

// GetFloatPS looks and get child float by given path and separator.
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

// GetIntPS looks and get child integer by given path and separator.
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

// GetUintPS looks and get child unsigned int by given path and separator.
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
