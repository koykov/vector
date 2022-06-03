package vector

// Shorthands for node.Get*PS() methods with "." used as separator.

// Dot looks and get child node by given path and "." separator.
func (n *Node) Dot(path string) *Node {
	return n.GetPS(path, ".")
}

// DotObject looks and get child object by given path and "." separator.
func (n *Node) DotObject(path string) *Node {
	return n.GetObjectPS(path, ".")
}

// DotArray looks and get child array by given path and "." separator.
func (n *Node) DotArray(path string) *Node {
	return n.GetArrayPS(path, ".")
}

// DotBytes looks and get child bytes by given path and "." separator.
func (n *Node) DotBytes(path string) []byte {
	return n.GetBytesPS(path, ".")
}

// DotString looks and get child string by given path and "." separator.
func (n *Node) DotString(path string) string {
	return n.GetStringPS(path, ".")
}

// DotBool looks and get child bool by given path and "." separator.
func (n *Node) DotBool(path string) bool {
	return n.GetBoolPS(path, ".")
}

// DotFloat looks and get child float by given path and "." separator.
func (n *Node) DotFloat(path string) (float64, error) {
	return n.GetFloatPS(path, ".")
}

// DotInt looks and get child integer by given path and "." separator.
func (n *Node) DotInt(path string) (int64, error) {
	return n.GetIntPS(path, ".")
}

// DotUint looks and get child unsigned int by given path and "." separator.
func (n *Node) DotUint(path string) (uint64, error) {
	return n.GetUintPS(path, ".")
}
