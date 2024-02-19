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
}
