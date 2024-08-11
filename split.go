package vector

import (
	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/entry"
)

// port of bytealg.AppendSplitEntryString, but with additional logic for checking square brackets and "@" separator.
func (vec *Vector) appendSplitPath(dst []entry.Entry64, s, sep string) []entry.Entry64 {
	if len(s) == 0 {
		return dst
	}
	off := len(vec.buf)
	vec.buf = append(vec.buf, s...)
	b := vec.buf[off:]

	off = 0
	for p := bytealg.IndexByteAtBytes(b, '[', off); p != -1; {
		b[p] = '.'
		off = p + 1
	}
	off = 0
	for p := bytealg.IndexByteAtBytes(b, ']', off); p != -1; {
		b[p] = '.'
		off = p + 1
	}
	off = 0
	for p := bytealg.IndexByteAtBytes(b, '@', off); p != -1; {
		b[p] = '.'
		off = p + 1
	}

	dst = bytealg.AppendSplitEntryBytes(dst, b, byteconv.S2B(sep), -1)

	for i := 0; i < len(dst); i++ {
		if lo, hi := dst[i].Decode(); lo == hi {
			dst = append(dst[:i], dst[i+1:]...)
		}
	}

	return dst
}
