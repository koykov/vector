package vector

import (
	"math"

	"github.com/koykov/bytealg"
	"github.com/koykov/byteconv"
	"github.com/koykov/entry"
)

// A wrapper around bytealg.AppendSplitEntryString with additional logic for checking square brackets and "@" separator.
func (vec *Vector) appendSplitPath(dst []entry.Entry64, s, sep string) (string, []entry.Entry64) {
	_ = splitTable[math.MaxUint8]
	n := len(s)
	if n == 0 {
		return s, dst
	}

	off := len(vec.buf)
	_ = s[n-1]
	for i := 0; i < n; i++ {
		if splitTable[s[i]] {
			vec.buf = append(vec.buf, sep...)
			continue
		}
		vec.buf = append(vec.buf, s[i])
	}

	b := vec.buf[off:]
	dst = bytealg.AppendSplitEntryBytes(dst, b, byteconv.S2B(sep), -1)

	for i := 0; i < len(dst); i++ {
		if lo, hi := dst[i].Decode(); lo == hi {
			dst = append(dst[:i], dst[i+1:]...)
		}
	}

	return byteconv.B2S(b), dst
}

var splitTable = [math.MaxUint8 + 1]bool{}

func init() {
	splitTable['@'] = true
	splitTable['['] = true
	splitTable[']'] = true
}
