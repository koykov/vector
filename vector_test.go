package vector

import (
	"fmt"
	"sync"
	"testing"
)

var testPool = sync.Pool{New: func() any { return &Vector{} }}

func TestVector(t *testing.T) {
	t.Run("populate", func(t *testing.T) {
		var vec Vector
		_ = vec.SetSrc([]byte("zzz"), false)
		r, ri := vec.AcquireNode(0)
		r.SetType(TypeObject)

		n, i := vec.AcquireChildWithType(r, 1, TypeString)
		n.Key().InitString("foo", 0, 3)
		n.Value().InitString("qwerty", 0, 6)
		r.ReleaseChild(i, n)

		n, i = vec.AcquireChildWithType(r, 1, TypeString)
		n.Key().InitString("bar", 0, 3)
		n.Value().InitString("asdfgh", 0, 6)
		r.ReleaseChild(i, n)

		vec.ReleaseNode(ri, r)

		m := make(map[string]any)
		vec.Populate(m)

		for k, v := range m {
			s, ok := v.(fmt.Stringer)
			if !ok {
				t.FailNow()
			}
			switch k {
			case "foo":
				if s.String() != "qwerty" {
					t.FailNow()
				}
			case "bar":
				if s.String() != "asdfgh" {
					t.FailNow()
				}
			}
		}
	})
}
