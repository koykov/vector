package vector

// Shorthands for node.Get*PS() methods with "." used as separator.

// Look and get child node by given path and "." separator.
func (n *Node) Dot(path string) *Node {
	return n.GetPS(path, ".")
}

// Look and get child object by given path and "." separator.
func (n *Node) DotObject(path string) *Node {
	return n.GetObjectPS(path, ".")
}

// Look and get child array by given path and "." separator.
func (n *Node) DotArray(path string) *Node {
	return n.GetArrayPS(path, ".")
}

// Look and get child bytes by given path and "." separator.
func (n *Node) DotBytes(path string) []byte {
	return n.GetBytesPS(path, ".")
}

// Look and get child string by given path and "." separator.
func (n *Node) DotString(path string) string {
	return n.GetStringPS(path, ".")
}

// Look and get child bool by given path and "." separator.
func (n *Node) DotBool(path string) bool {
	return n.GetBoolPS(path, ".")
}

// Look and get child float by given path and "." separator.
func (n *Node) DotFloat(path string) (float64, error) {
	return n.GetFloatPS(path, ".")
}

// Look and get child integer by given path and "." separator.
func (n *Node) DotInt(path string) (int64, error) {
	return n.GetIntPS(path, ".")
}

// Look and get child unsigned int by given path and "." separator.
func (n *Node) DotUint(path string) (uint64, error) {
	return n.GetUintPS(path, ".")
}
