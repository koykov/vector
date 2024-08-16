package vector

import (
	"strconv"
)

// BufLen returns length of buffer.
func (vec *Vector) BufLen() int {
	return len(vec.buf)
}

// Buf returns raw buffer bytes.
func (vec *Vector) Buf() []byte {
	return vec.buf
}

// BufReplace replaces buffer with b.
func (vec *Vector) BufReplace(b []byte) {
	vec.buf = b
}

// BufUpdateWith replaces buffer with b.
// DEPRECATED: use BufReplace instead.
func (vec *Vector) BufUpdateWith(b []byte) {
	vec.buf = b
}

// Bufferize appends b to internal buffer and returns buffered value.
func (vec *Vector) Bufferize(b []byte) []byte {
	off := vec.BufLen()
	vec.buf = append(vec.buf, b...)
	return vec.buf[off:]
}

// BufferizeString appends string to internal buffer and returns buffered value.
func (vec *Vector) BufferizeString(s string) []byte {
	off := vec.BufLen()
	vec.buf = append(vec.buf, s...)
	return vec.buf[off:]
}

// BufferizeByte appends b to internal buffer and returns buffered value.
func (vec *Vector) BufferizeByte(b byte) []byte {
	off := vec.BufLen()
	vec.buf = append(vec.buf, b)
	return vec.buf[off:]
}

// BufferizeInt appends integer to internal buffer and returns buffered value.
func (vec *Vector) BufferizeInt(i int64) []byte {
	off := vec.BufLen()
	vec.buf = strconv.AppendInt(vec.buf, i, 10)
	return vec.buf[off:]
}

// BufferizeUint appends unsigned integer to internal buffer and returns buffered value.
func (vec *Vector) BufferizeUint(u uint64) []byte {
	off := vec.BufLen()
	vec.buf = strconv.AppendUint(vec.buf, u, 10)
	return vec.buf[off:]
}

// BufferizeFloat appends unsigned integer to internal buffer and returns buffered value.
func (vec *Vector) BufferizeFloat(f float64) []byte {
	off := vec.BufLen()
	vec.buf = strconv.AppendFloat(vec.buf, f, 'f', -1, 64)
	return vec.buf[off:]
}

// BufferizeFloatTune appends float to internal buffer and returns buffered value.
func (vec *Vector) BufferizeFloatTune(f float64, fmt byte, prec, bitSize int) []byte {
	off := vec.BufLen()
	vec.buf = strconv.AppendFloat(vec.buf, f, fmt, prec, bitSize)
	return vec.buf[off:]
}
