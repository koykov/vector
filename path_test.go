package vector

import (
	"reflect"
	"testing"

	"github.com/koykov/entry"
)

type pathStage struct {
	path   string
	expect []entry.Entry64
}

var pathStages = []pathStage{
	{path: "foobar", expect: []entry.Entry64{6}},
	{path: "foo.bar", expect: []entry.Entry64{3, 17179869191}},
	{path: "@version", expect: []entry.Entry64{8}},
	{path: "root@version", expect: []entry.Entry64{4, 17179869196}},
	{path: "root.qwe.rty@version", expect: []entry.Entry64{4, 21474836488, 38654705676, 51539607572}},
	{path: "foobar@", expect: []entry.Entry64{6}},
	{path: "foo.bar[2]", expect: []entry.Entry64{3, 17179869191, 34359738377}},
	{path: "foo[2].bar", expect: []entry.Entry64{3, 17179869189, 30064771082}},
	{path: "foo.bar[15].baz@qwe", expect: []entry.Entry64{3, 17179869191, 34359738378, 51539607567, 64424509459}},
}

func TestPath(t *testing.T) {
	for _, stg := range pathStages {
		t.Run(stg.path, func(t *testing.T) {
			vec := Vector{}
			vec.splitPath(stg.path, ".")
			if !reflect.DeepEqual(stg.expect, vec.bufKE) {
				t.FailNow()
			}
		})
	}
}

func BenchmarkPath(b *testing.B) {
	for _, stg := range pathStages {
		b.Run(stg.path, func(b *testing.B) {
			b.ReportAllocs()
			vec := Vector{}
			for i := 0; i < b.N; i++ {
				vec.splitPath(stg.path, ".")
			}
		})
	}
}
