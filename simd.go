package vector

import "github.com/koykov/simd/indexbyte"

// IndexByte returns position of b in s. Return -1 if not found.
func IndexByte(s []byte, b byte) int {
	return indexbyte.Index(s, b)
}

// IndexByteAt returns position of b in s starting from position at. Return -1 if not found.
func IndexByteAt(s []byte, b byte, at int) int {
	return indexbyte.IndexAt(s, b, at)
}

// IndexByteNE returns position of non-escaped b in s. Return -1 if not found.
func IndexByteNE(s []byte, b byte) int {
	return indexbyte.IndexNE(s, b)
}

// IndexByteAtNE returns position of non-escaped b in s from position at. Return -1 if not found.
func IndexByteAtNE(s []byte, b byte, at int) int {
	return indexbyte.IndexAtNE(s, b, at)
}

var _, _, _, _ = IndexByte, IndexByteAt, IndexByteNE, IndexByteAtNE
