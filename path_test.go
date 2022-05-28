package vector

import (
	"reflect"
	"testing"
)

func TestPath(t *testing.T) {
	type stage struct {
		path   string
		expect []string
	}
	stages := []stage{
		{path: "foobar", expect: []string{"foobar"}},
		{path: "foo.bar", expect: []string{"foo", "bar"}},
		{path: "@version", expect: []string{"version"}},
		{path: "root@version", expect: []string{"root", "@version"}},
		{path: "root.qwe.rty@version", expect: []string{"root", "qwe", "rty", "@version"}},
		{path: "foobar@", expect: []string{"foobar"}},
	}
	for _, stg := range stages {
		t.Run(stg.path, func(t *testing.T) {
			vec := Vector{}
			vec.splitPath(stg.path, ".")
			if !reflect.DeepEqual(stg.expect, vec.bufSS) {
				t.FailNow()
			}
		})
	}
}
