package vector

import "bytes"

// EqualWith compares vector with exp.
func (vec *Vector) EqualWith(exp *Vector) bool {
	return equal(vec.Root(), exp.Root())
}

// EqualWith compares node with exp.
func (n *Node) EqualWith(exp *Node) bool {
	return equal(n, exp)
}

func equal(a, b *Node) bool {
	switch {
	case a == nil && b == nil:
		return true
	case a == nil:
		return false
	case b == nil:
		return false
	}
	if a.Type() != b.Type() {
		return false
	}
	if a.Limit() != b.Limit() {
		return false
	}
	ok := true
	switch a.Type() {
	case TypeObject:
		a.Each(func(_ int, ac *Node) {
			bc := b.Get(ac.KeyString())
			if bc.Type() == TypeNull {
				ok = false
				return
			}
			if !equal(ac, bc) {
				ok = false
				return
			}
		})
	case TypeArray:
		a.Each(func(idx int, ac *Node) {
			bc := b.At(idx)
			if !equal(ac, bc) {
				ok = false
				return
			}
		})
	case TypeUnknown:
	case TypeNull:
		return true
	default:
		return bytes.Equal(a.Value().RawBytes(), b.Value().RawBytes())
	}
	return ok
}
