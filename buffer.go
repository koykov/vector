package vector

import (
	"strconv"

	"github.com/koykov/byteconv"
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
func (vec *Vector) BufferizeString(s string) string {
	off := vec.BufLen()
	vec.buf = append(vec.buf, s...)
	return byteconv.B2S(vec.buf[off:])
}

// BufAppend appends bytes to the buffer.
// DEPRECATED: use Bufferize instead.
func (vec *Vector) BufAppend(s []byte) {
	vec.buf = append(vec.buf, s...)
}

// BufAppendStr appends string to the buffer.
func (vec *Vector) BufAppendStr(s string) {
	vec.buf = append(vec.buf, s...)
}

// BufAppendByte appends single byte to the buffer.
func (vec *Vector) BufAppendByte(b byte) {
	vec.buf = append(vec.buf, b)
}

// BufAppendInt appends int to the buffer.
func (vec *Vector) BufAppendInt(i int64) {
	vec.buf = strconv.AppendInt(vec.buf, i, 10)
}

// BufAppendUint appends uint to the buffer.
func (vec *Vector) BufAppendUint(u uint64) {
	vec.buf = strconv.AppendUint(vec.buf, u, 10)
}

// BufAppendFloat appends float to the buffer.
func (vec *Vector) BufAppendFloat(f float64) {
	vec.buf = strconv.AppendFloat(vec.buf, f, 'f', -1, 64)
}

// BufAppendFloatTune appends float with extended params to the buffer.
func (vec *Vector) BufAppendFloatTune(f float64, fmt byte, prec, bitSize int) {
	vec.buf = strconv.AppendFloat(vec.buf, f, fmt, prec, bitSize)
}
