package vector

type Vector struct {
	buf, src []byte
	bufSS    []string
	addr     uint64
	selfPtr  uintptr
	nodes    []Node
	nodeL    int
	errLn    int
	index    index
}
