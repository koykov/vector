package vector

import "strconv"

// Get length of buffer.
func (vec *Vector) BufLen() int {
	return len(vec.buf)
}

// Get raw buffer bytes.
func (vec *Vector) Buf() []byte {
	return vec.buf
}

// Append bytes to the buffer.
func (vec *Vector) BufAppend(s []byte) {
	vec.buf = append(vec.buf, s...)
}

// Append string to the buffer.
func (vec *Vector) BufAppendStr(s string) {
	vec.buf = append(vec.buf, s...)
}

// Append int to the buffer.
func (vec *Vector) BufAppendInt(i int64) {
	vec.buf = strconv.AppendInt(vec.buf, i, 10)
}

// Append uint to the buffer.
func (vec *Vector) BufAppendUint(u uint64) {
	vec.buf = strconv.AppendUint(vec.buf, u, 10)
}

// Append float to the buffer.
func (vec *Vector) BufAppendFloat(f float64) {
	vec.buf = strconv.AppendFloat(vec.buf, f, 'f', -1, 64)
}
