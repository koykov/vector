package vector

import (
	"testing"
	"unsafe"
)

func TestByteptr(t *testing.T) {
	t.Run("size", func(t *testing.T) {
		var p Byteptr
		if unsafe.Sizeof(p) != byteptrSize {
			t.FailNow()
		}
	})
}
