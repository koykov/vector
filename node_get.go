package vector

import (
	"strconv"

	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
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
	node := n
	if n.typ == TypeAlias {
		if idxs := n.ChildrenIndices(); len(idxs) > 0 {
			node = &vec.nodes[0]
		}
	}
	if node.typ == TypeObj {
		for i := node.offset; i < node.limit; i++ {
			idx := vec.Index.val(node.depth+1, i)
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
	if node.typ == TypeArr {
		i, err := strconv.Atoi(keys[0])
		if err != nil || i >= node.limit {
			return nullNode
		}
		idx := vec.Index.val(node.depth+1, node.offset+i)
		child := &vec.nodes[idx]
		tail := keys[1:]
		if len(tail) == 0 {
			return child
		} else if n.val.Len() > 0 && len(tail) == 1 && node.val.String() == tail[0] {
			return node
		} else {
			return child.Get(tail...)
		}
	}
	return nullNode
}

// Entry based version of Get.
func (n *Node) getKE(path string, keys ...entry.Entry64) *Node {
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
	node := n
	if n.typ == TypeAlias {
		if idxs := n.ChildrenIndices(); len(idxs) > 0 {
			node = &vec.nodes[0]
		}
	}
	if node.typ == TypeObj {
		for i := node.offset; i < node.limit; i++ {
			idx := vec.Index.val(node.depth+1, i)
			child := &vec.nodes[idx]
			if child.keyEqualKE(path, keys[0]) {
				tail := keys[1:]
				if len(tail) == 0 {
					return child
				} else {
					return child.getKE(path, tail...)
				}
			}
		}
	}
	if node.typ == TypeArr {
		lo, hi := keys[0].Decode()
		skey := path[lo:hi]
		i, err := strconv.Atoi(skey)
		if err != nil || i >= node.limit {
			return nullNode
		}
		idx := vec.Index.val(node.depth+1, node.offset+i)
		child := &vec.nodes[idx]
		tail := keys[1:]
		if len(tail) == 0 {
			return child
		} else if node.val.Len() > 0 && len(tail) == 1 && node.val.equalKE(path, tail[0]) {
			return node
		} else {
			return child.getKE(path, tail...)
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
	vec.bufKE = bytealg.AppendSplitEntryString(vec.bufKE[:0], path, separator, -1)
	return n.getKE(path, vec.bufKE...)
}

// GetObjectPS looks and get child object by given path and separator.
func (n *Node) GetObjectPS(path, separator string) *Node {
	vec := n.indirectVector()
	if vec == nil {
		return nullNode
	}
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
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
	path = vec.splitPath(path, separator)
	node := n.getKE(path, vec.bufKE...)
	if node.typ == TypeNull {
		return 0, ErrNotFound
	}
	if node.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return node.Uint()
}
