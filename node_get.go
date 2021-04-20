package vector

import "github.com/koykov/bytealg"

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
