package vector

// Shorthands for vector.Get*PS() methods with "." used as separator.

// Dot looks and get node by given path and "." separator.
func (vec *Vector) Dot(path string) *Node {
	return vec.GetPS(path, ".")
}

// DotObject looks and get object by given path and "." separator.
func (vec *Vector) DotObject(path string) *Node {
	return vec.GetObjectPS(path, ".")
}

// DotArray looks and get array by given path and "." separator.
func (vec *Vector) DotArray(path string) *Node {
	return vec.GetArrayPS(path, ".")
}

// DotBytes looks and get bytes by given path and "." separator.
func (vec *Vector) DotBytes(path string) []byte {
	return vec.GetBytesPS(path, ".")
}

// DotString looks and get string by given path and "." separator.
func (vec *Vector) DotString(path string) string {
	return vec.GetStringPS(path, ".")
}

// DotBool looks and get bool by given path and "." separator.
func (vec *Vector) DotBool(path string) bool {
	return vec.GetBoolPS(path, ".")
}

// DotFloat looks and get float by given path and "." separator.
func (vec *Vector) DotFloat(path string) (float64, error) {
	return vec.GetFloatPS(path, ".")
}

// DotInt looks and get integer by given path and "." separator.
func (vec *Vector) DotInt(path string) (int64, error) {
	return vec.GetIntPS(path, ".")
}

// DotUint looks and get unsigned integer by given path and "." separator.
func (vec *Vector) DotUint(path string) (uint64, error) {
	return vec.GetUintPS(path, ".")
}
