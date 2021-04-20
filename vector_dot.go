package vector

func (vec *Vector) DotObject(path string) *Node {
	return vec.GetObjectPS(path, ".")
}

func (vec *Vector) DotArray(path string) *Node {
	return vec.GetArrayPS(path, ".")
}

func (vec *Vector) DotBytes(path string) []byte {
	return vec.GetBytesPS(path, ".")
}

func (vec *Vector) DotString(path string) string {
	return vec.GetStringPS(path, ".")
}

func (vec *Vector) DotBool(path string) bool {
	return vec.GetBoolPS(path, ".")
}

func (vec *Vector) DotFloat(path string) (float64, error) {
	return vec.GetFloatPS(path, ".")
}

func (vec *Vector) DotInt(path string) (int64, error) {
	return vec.GetIntPS(path, ".")
}

func (vec *Vector) DotUint(path string) (uint64, error) {
	return vec.GetUintPS(path, ".")
}
