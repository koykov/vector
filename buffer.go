package vector

import "strconv"

// BufLen returns length of buffer.
func (vec *Vector) BufLen() int {
	return len(vec.buf)
}

// Buf returns raw buffer bytes.
func (vec *Vector) Buf() []byte {
	return vec.buf
}

// BufUpdateWith replaces buffer with b.
func (vec *Vector) BufUpdateWith(b []byte) {
	vec.buf = b
}

// BufAppend appends bytes to the buffer.
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
