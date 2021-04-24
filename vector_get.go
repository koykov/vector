package vector

import (
	"strconv"

	"github.com/koykov/bytealg"
)

func (vec *Vector) Get(keys ...string) *Node {
	if len(keys) == 0 {
		if vec.Len() > 0 {
			return vec.nodes[0]
		}
		return nullNode
	}

	r := vec.nodes[0]
	if r.typ != TypeObj && r.typ != TypeArr {
		if len(keys) > 1 {
			return nullNode
		}
		return r
	}

	if r.typ == TypeArr {
		return vec.getArr(r, keys...)
	}
	if r.typ == TypeObj {
		return vec.getObj(r, keys...)
	}
	return r
}

func (vec *Vector) GetObject(keys ...string) *Node {
	v := vec.Get(keys...)
	if v.Type() != TypeObj {
		return nullNode
	}
	return v.Object()
}

func (vec *Vector) GetArray(keys ...string) *Node {
	v := vec.Get(keys...)
	if v.Type() != TypeArr {
		return nullNode
	}
	return v.Array()
}

func (vec *Vector) GetBytes(keys ...string) []byte {
	v := vec.Get(keys...)
	if v.Type() != TypeStr {
		return nil
	}
	return v.Bytes()
}

func (vec *Vector) GetString(keys ...string) string {
	v := vec.Get(keys...)
	if v.Type() != TypeStr {
		return ""
	}
	return v.String()
}

func (vec *Vector) GetBool(keys ...string) bool {
	v := vec.Get(keys...)
	if v.Type() != TypeBool {
		return false
	}
	return v.Bool()
}

func (vec *Vector) GetFloat(keys ...string) (float64, error) {
	v := vec.Get(keys...)
	if v.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Float()
}

func (vec *Vector) GetInt(keys ...string) (int64, error) {
	v := vec.Get(keys...)
	if v.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Int()
}

func (vec *Vector) GetUint(keys ...string) (uint64, error) {
	v := vec.Get(keys...)
	if v.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Uint()
}

func (vec *Vector) GetPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	return vec.Get(vec.bufSS...)
}

func (vec *Vector) GetObjectPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() != TypeObj {
		return nullNode
	}
	return v.Object()
}

func (vec *Vector) GetArrayPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() != TypeArr {
		return nullNode
	}
	return v.Array()
}

func (vec *Vector) GetBytesPS(path, separator string) []byte {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() != TypeStr {
		return nil
	}
	return v.Bytes()
}

func (vec *Vector) GetStringPS(path, separator string) string {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() != TypeStr {
		return ""
	}
	return v.String()
}

func (vec *Vector) GetBoolPS(path, separator string) bool {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() != TypeBool {
		return false
	}
	return v.Bool()
}

func (vec *Vector) GetFloatPS(path, separator string) (float64, error) {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Float()
}

func (vec *Vector) GetIntPS(path, separator string) (int64, error) {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Int()
}

func (vec *Vector) GetUintPS(path, separator string) (uint64, error) {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	v := vec.Get(vec.bufSS...)
	if v.Type() == TypeUnk {
		return 0, ErrNotFound
	}
	if v.Type() != TypeNum {
		return 0, ErrIncompatType
	}
	return v.Uint()
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
	v := vec.nodes[i]
	tail := keys[1:]
	if v.typ != TypeArr && v.typ != TypeObj {
		if len(tail) > 0 {
			return nullNode
		}
		return v
	}
	if v.typ == TypeArr {
		return vec.getArr(v, tail...)
	}
	if v.typ == TypeObj {
		return vec.getObj(v, tail...)
	}
	return nil
}

func (vec *Vector) getObj(root *Node, keys ...string) *Node {
	if len(keys) == 0 {
		return root
	}
	var v *Node
	for i := root.offset; i < root.limit; i++ {
		k := vec.Index.val(root.depth+1, i)
		v = vec.nodes[k]
		if v.key.String() == keys[0] {
			break
		}
	}
	if v == nil {
		return v
	}
	tail := keys[1:]
	if v.typ != TypeArr && v.typ != TypeObj {
		if len(tail) > 0 {
			return nullNode
		}
		return v
	}
	if v.typ == TypeArr {
		return vec.getArr(v, tail...)
	}
	if v.typ == TypeObj {
		return vec.getObj(v, tail...)
	}
	return nil
}
