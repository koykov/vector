package vector

// Shorthands for vector.Get*PS() methods with "." used as separator.

// Look and get node by given path and "." separator.
func (vec *Vector) Dot(path string) *Node {
	return vec.GetPS(path, ".")
}

// Look and get object by given path and "." separator.
func (vec *Vector) DotObject(path string) *Node {
	return vec.GetObjectPS(path, ".")
}

// Look and get array by given path and "." separator.
func (vec *Vector) DotArray(path string) *Node {
	return vec.GetArrayPS(path, ".")
}

// Look and get bytes by given path and "." separator.
func (vec *Vector) DotBytes(path string) []byte {
	return vec.GetBytesPS(path, ".")
}

// Look and get string by given path and "." separator.
func (vec *Vector) DotString(path string) string {
	return vec.GetStringPS(path, ".")
}

// Look and get bool by given path and "." separator.
func (vec *Vector) DotBool(path string) bool {
	return vec.GetBoolPS(path, ".")
}

// Look and get float by given path and "." separator.
func (vec *Vector) DotFloat(path string) (float64, error) {
	return vec.GetFloatPS(path, ".")
}

// Look and get integer by given path and "." separator.
func (vec *Vector) DotInt(path string) (int64, error) {
	return vec.GetIntPS(path, ".")
}

// Look and get unsigned integer by given path and "." separator.
func (vec *Vector) DotUint(path string) (uint64, error) {
	return vec.GetUintPS(path, ".")
}
