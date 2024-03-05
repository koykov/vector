package vector

import (
	"testing"
	"unsafe"
)

func TestNode(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		var n Node
		if sz := unsafe.Sizeof(n); sz != nodeSize {
			t.Errorf("node size fail: need %d, got %d", nodeSize, sz)
		}
	})
}
