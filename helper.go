package vector

import "io"

// Helper object interface.
type Helper interface {
	// Indirect convert byteptr to byte slice and apply custom logic.
	Indirect(*Byteptr) []byte
	// Beautify makes a beauty view of node.
	Beautify(io.Writer, *Node) error
	// Marshal serializes node.
	Marshal(io.Writer, *Node) error
}

// SetHelper sets helper to vector object.
func (vec *Vector) SetHelper(helper Helper) {
	vec.Helper = helper
}
