package vector

import (
	"strings"

	"github.com/koykov/bytealg"
	"github.com/koykov/entry"
)

// port of bytealg.AppendSplitEntryString, but with additional logic for checking square brackets and "@" separator.
func appendSplitPath(buf []entry.Entry64, s, sep string) []entry.Entry64 {
	if len(s) == 0 {
		return buf
	}
	var off int
	var m int
	for m < len(s) {
		m = bytealg.IndexAtString(s, sep, off)
		if m < 0 {
			m = len(s)
		}
		k := s[:m]
		if p := strings.IndexByte(k, '['); p != -1 {
			// todo check square brackets
		} else if p := strings.IndexByte(k, '@'); p > 0 {
			// todo check @ symbol
		} else {
			var e entry.Entry64
			e.Encode(uint32(off), uint32(off+m))
			buf = append(buf, e)
		}
		s = s[m+len(sep):]
		off += m + len(sep)
	}
	return buf
}
