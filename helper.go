package vector

// Helper object interface.
type Helper interface {
	// Indirect convert byteptr to byte slice and apply custom logic.
	Indirect(*Byteptr) []byte
}

// SetHelper sets helper to vector object.
func (vec *Vector) SetHelper(helper Helper) {
	vec.Helper = helper
}
