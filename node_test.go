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
		vec := testPool.Get().(*Vector)
		defer testPool.Put(vec)

		_ = vec.SetSrc([]byte("N/D"), false) // emulate parsing to init vector
		root, ri := vec.AcquireNodeWithType(0, TypeObject)

		sn, si := vec.AcquireChildWithType(root, 1, TypeString)
		sn.Key().InitString("foo", 0, 3)
		sn.Value().InitString("bar", 0, 3)
		vec.ReleaseNode(si, sn)

		an, ai := vec.AcquireChildWithType(root, 1, TypeAlias)
		an.Key().InitString("qwe", 0, 3)
		an.AliasOf(sn)
		vec.ReleaseNode(ai, an)

		vec.ReleaseNode(ri, root)

		val := vec.Dot("qwe").String()
		if val != "bar" {
			t.FailNow()
		}
	})
}
