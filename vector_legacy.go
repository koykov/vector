package vector

import "unsafe"

// ParseStr is a legacy version of ParseString.
// DEPRECATED: use ParseString instead.
func (vec *Vector) ParseStr(_ string) error {
	return ErrNotImplement
}

// ParseCopyStr is a legacy version of ParseCopyString.
// DEPRECATED: use ParseCopyString instead.
func (vec *Vector) ParseCopyStr(_ string) error {
	return ErrNotImplement
}

// GetNode is a legacy version of AcquireNode.
// DEPRECATED: use AcquireNode instead.
func (vec *Vector) GetNode(depth int) (node *Node, idx int) {
	node, idx = vec.ackNode(depth)
	vec.Index.Register(depth, idx)
	return
}

// GetNodeWT is a legacy version of AcquireNodeWithType.
// DEPRECATED: use AcquireNodeWithType instead.
func (vec *Vector) GetNodeWT(depth int, typ Type) (*Node, int) {
	node, idx := vec.AcquireNode(depth)
	node.typ = typ
	return node, idx
}

// GetChild is a legacy version of AcquireChild.
// DEPRECATED: use AcquireChild instead.
func (vec *Vector) GetChild(root *Node, depth int) (*Node, int) {
	return vec.AcquireChildWithType(root, depth, TypeUnknown)
}

// GetChildWT is a legacy version of AcquireChildWithType.
// DEPRECATED: use AcquireChildWithType instead.
func (vec *Vector) GetChildWT(root *Node, depth int, typ Type) (*Node, int) {
	node, idx := vec.ackNode(depth)
	node.typ = typ
	node.pptr = root.ptr()
	root.SetLimit(vec.Index.Register(depth, idx))
	return node, idx
}

// PutNode is a legacy version of ReleaseNode.
// DEPRECATED: use ReleaseNode instead.
func (vec *Vector) PutNode(idx int, node *Node) {
	l := unsafe.Pointer(&vec.nodes[idx])
	r := unsafe.Pointer(node)
	if uintptr(l) != uintptr(r) {
		vec.nodes[idx] = *node
	}
}
