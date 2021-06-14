package vector

// Helper object interface.
type Helper interface {
	// Convert byteptr to byte slice and apply custom logic.
	ConvertByteptr(*Byteptr) []byte
}

// Assign helper to vector object.
func (vec *Vector) SetHelper(helper Helper) {
	vec.Helper = helper
}
