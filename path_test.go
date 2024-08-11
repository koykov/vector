package vector

import (
	"reflect"
	"testing"

	"github.com/koykov/entry"
)

func TestPath(t *testing.T) {
	type stage struct {
		path   string
		expect []entry.Entry64
	}

	stages := []stage{
		// {path: "foobar", expect: []entry.Entry64{6}},
		// {path: "foo.bar", expect: []entry.Entry64{3, 17179869191}},
		{path: "@version", expect: []entry.Entry64{8}},
		// {path: "root@version", expect: []entry.Entry64{4, 17179869196}},
		// {path: "root.qwe.rty@version", expect: []entry.Entry64{4, 21474836488, 38654705676, 51539607572}},
		// {path: "foobar@", expect: []entry.Entry64{6}},
		// {path: "foo.bar[2]", expect: []entry.Entry64{3, 17179869191, 17179869191}},
		// {path: "foo[2].bar", expect: []entry.Entry64{3, 17179869191, 17179869191}},
	}
	for _, stg := range stages {
		t.Run(stg.path, func(t *testing.T) {
			vec := Vector{}
			vec.splitPath(stg.path, ".")
			if !reflect.DeepEqual(stg.expect, vec.bufKE) {
				t.FailNow()
			}
		})
	}
}
