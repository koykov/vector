package vector

import (
	"unsafe"

	"github.com/koykov/bytealg"
)

type Vector struct {
	buf, src []byte
	bufSS    []string
	addr     uint64
	selfPtr  uintptr
	nodes    []Node
	nodeL    int
	errOff   int
	index    index
}

func (vec *Vector) Len() int {
	return vec.nodeL
}

func (vec *Vector) ErrorOffset() int {
	return vec.errOff
}

func (vec *Vector) Root() *Node {
	return vec.Get()
}

func (vec *Vector) Get(keys ...string) *Node {
	if len(keys) == 0 {
		if vec.Len() > 0 {
			return &vec.nodes[0]
		}
		return nullNode
	}

	r := &vec.nodes[0]
	if r.typ != TypeObj && r.typ != TypeArr {
		if len(keys) > 1 {
			return nil
		}
		return r
	}

	if r.typ == TypeArr {
		return vec.getArr(r, keys...)
	}
	if r.typ == TypeObj {
		return vec.getObj(r, keys...)
	}
	return r
}

func (vec *Vector) GetPS(path, separator string) *Node {
	vec.bufSS = bytealg.AppendSplitStr(vec.bufSS[:0], path, separator, -1)
	return vec.Get(vec.bufSS...)
}

func (vec *Vector) Dot(path string) *Node {
	return vec.GetPS(path, ".")
}

func (vec *Vector) Exists(key string) bool {
	n := vec.Root()
	return n.Exists(key)
}

func (vec *Vector) AcquireNode(depth int) (r *Node) {
	if vec.nodeL < len(vec.nodes) {
		r = &vec.nodes[vec.nodeL]
		r.Reset()
		vec.nodeL++
	} else {
		r = &Node{typ: TypeUnk}
		vec.nodes = append(vec.nodes, *r)
		vec.nodeL++
	}
	r.vecPtr, r.depth = vec.ptr(), depth
	return
}

func (vec *Vector) Reset() {
	if vec.nodeL == 0 {
		return
	}
	_ = vec.nodes[vec.nodeL-1]
	for i := 0; i < vec.nodeL; i++ {
		vec.nodes[i].vecPtr = 0
	}
	vec.buf, vec.src = vec.buf[:0], nil
	vec.addr, vec.nodeL, vec.errOff = 0, 0, 0
	vec.index.reset()
}

func (vec *Vector) ptr() uintptr {
	if vec.selfPtr == 0 {
		vec.selfPtr = uintptr(unsafe.Pointer(vec))
	}
	return vec.selfPtr
}
