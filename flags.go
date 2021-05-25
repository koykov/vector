package vector

// Flags bitmask type.
type Flags uint8

func (f *Flags) SetFlag(flag uint8, value bool) {
	if value {
		*f = Flags((uint8)(*f) | flag)
	} else {
		*f = Flags((uint8)(*f) & ^flag)
	}
}

func (f *Flags) CheckFlag(flag uint8) bool {
	return (uint8)(*f)&flag == 1
}

func (f *Flags) Reset() {
	*f = 0
}
