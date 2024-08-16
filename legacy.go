package vector

import "strconv"

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

// ParseStr is a legacy version of ParseString.
// DEPRECATED: use ParseString instead.
func (vec *Vector) ParseStr(_ string) error {
	return ErrNotImplement
}

// ParseCopyStr is a legacy version of ParseCopyString.
// DEPRECATED: use ParseCopyString instead.
func (vec *Vector) ParseCopyStr(_ string) error {
	return ErrNotImplement
}

// BufUpdateWith is a legacy version of BufReplaceWith.
// DEPRECATED: use BufReplaceWith instead.
func (vec *Vector) BufUpdateWith(b []byte) {
	vec.buf = b
}

// BufAppend is a legacy version of Bufferize.
// DEPRECATED: use Bufferize instead.
func (vec *Vector) BufAppend(s []byte) {
	vec.buf = append(vec.buf, s...)
}

// BufAppendStr is a legacy version of BufferizeString.
// DEPRECATED: use BufferizeString instead.
func (vec *Vector) BufAppendStr(s string) {
	vec.buf = append(vec.buf, s...)
}

// BufAppendByte is a legacy version of BufferizeByte.
// DEPRECATED: use BufferizeByte instead.
func (vec *Vector) BufAppendByte(b byte) {
	vec.buf = append(vec.buf, b)
}

// BufAppendInt is a legacy version of BufferizeInt.
// DEPRECATED: use BufferizeInt instead.
func (vec *Vector) BufAppendInt(i int64) {
	vec.buf = strconv.AppendInt(vec.buf, i, 10)
}

// BufAppendUint is a legacy version of BufferizeUint.
// DEPRECATED: use BufferizeUint instead.
func (vec *Vector) BufAppendUint(u uint64) {
	vec.buf = strconv.AppendUint(vec.buf, u, 10)
}

// BufAppendFloat is a legacy version of BufferizeFloat.
// DEPRECATED: use BufferizeFloat instead.
func (vec *Vector) BufAppendFloat(f float64) {
	vec.buf = strconv.AppendFloat(vec.buf, f, 'f', -1, 64)
}

// BufAppendFloatTune is a legacy version of BufferizeFloatTune.
// DEPRECATED: use BufferizeFloatTune instead.
func (vec *Vector) BufAppendFloatTune(f float64, fmt byte, prec, bitSize int) {
	vec.buf = strconv.AppendFloat(vec.buf, f, fmt, prec, bitSize)
}
