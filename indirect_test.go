package vector

import "testing"

func TestIndirect(t *testing.T) {
	t.Run("vector", func(t *testing.T) {
		vec := Vector{}
		s := "foobar"
		root, i := vec.GetNodeWT(0, TypeStr)
		root.Value().InitStr(s, 0, len(s))
		vec.PutNode(i, root)
		ptr := vec.ptr()
		vec1 := indirectVector1(ptr)
		if vec1.Root().Value().String() != s {
			t.FailNow()
		}
	})
	t.Run("node", func(t *testing.T) {
		vec := Vector{}
		s := "foobar"
		root, i := vec.GetNodeWT(0, TypeStr)
		root.Value().InitStr(s, 0, len(s))
		vec.PutNode(i, root)
		ptr := vec.Root().ptr()
		n := indirectNode1(ptr)
		if n.Value().String() != s {
			t.FailNow()
		}
	})
}
