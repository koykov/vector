package vector

// Legacy types.
const (
	// TypeUnk is a legacy version of TypeUnknown.
	// DEPRECATED: use TypeUnknown instead.
	TypeUnk Type = 0
	// TypeObj is a legacy version of TypeObject.
	// DEPRECATED: use TypeObject instead.
	TypeObj Type = 1
	// TypeArr is a legacy version of TypeArray.
	// DEPRECATED: use TypeArray instead.
	TypeArr Type = 3
	// TypeStr is a legacy version of TypeString.
	// DEPRECATED: use TypeString instead.
	TypeStr Type = 4
	// TypeNum is a legacy version of TypeNumber.
	// DEPRECATED: use TypeNumber instead.
	TypeNum Type = 5
	// TypeAttr is a legacy version of TypeAttribute.
	// DEPRECATED: use TypeAttribute instead.
	TypeAttr Type = 7
)

// ParseStr parses source string.
// DEPRECATED: use ParseString instead.
func (vec *Vector) ParseStr(_ string) error {
	return ErrNotImplement
}

// ParseCopyStr copies source string and parse it.
// DEPRECATED: use ParseCopyString instead.
func (vec *Vector) ParseCopyStr(_ string) error {
	return ErrNotImplement
}
