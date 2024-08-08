package vector

import (
	"testing"
	"unsafe"
)

func TestNode(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		var n Node
		if unsafe.Sizeof(n) != nodeSize {
			t.FailNow()
		}
	})
	t.Run("alias", func(t *testing.T) {
		vec := &Vector{}
		_ = vec.SetSrc([]byte("N/D"), false) // emulate parsing to init vector
		root, ri := vec.GetNodeWT(0, TypeObj)

		sn, si := vec.GetChildWT(root, 1, TypeStr)
		sn.Key().InitString("foo", 0, 3)
		sn.Value().InitString("bar", 0, 3)
		vec.PutNode(si, sn)

		an, ai := vec.GetChildWT(root, 1, TypeAlias)
		an.Key().InitString("qwe", 0, 3)
		an.AliasOf(sn)
		vec.PutNode(ai, an)

		vec.PutNode(ri, root)

		val := vec.Dot("qwe").String()
		if val != "bar" {
			t.FailNow()
		}
	})
}
