package vector

func (n *Node) DotObject(path string) *Node {
	return n.GetObjectPS(path, ".")
}

func (n *Node) DotArray(path string) *Node {
	return n.GetArrayPS(path, ".")
}

func (n *Node) DotBytes(path string) []byte {
	return n.GetBytesPS(path, ".")
}

func (n *Node) DotString(path string) string {
	return n.GetStringPS(path, ".")
}

func (n *Node) DotBool(path string) bool {
	return n.GetBoolPS(path, ".")
}

func (n *Node) DotFloat(path string) (float64, error) {
	return n.GetFloatPS(path, ".")
}

func (n *Node) DotInt(path string) (int64, error) {
	return n.GetIntPS(path, ".")
}

func (n *Node) DotUint(path string) (uint64, error) {
	return n.GetUintPS(path, ".")
}
