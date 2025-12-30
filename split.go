package vector

import (
	"math"
	"strings"

	"github.com/koykov/entry"
	"github.com/koykov/simd/indextoken"
)

const splitPathThreshold = 128

// Split path by given separator.
//
// Caution! Don't use "@" as a separator, it will break work with attributes.
// TODO: consider escaped at symbol "\@".
func (vec *Vector) splitPath(path, separator string) {
	if len(separator) == 1 && separator[0] == '.' && len(path) > splitPathThreshold {
		vec.bufKE = vec.appendSplitPath(vec.bufKE[:0], path, separator)
		return
	}
	vec.bufKE = vec.appendSplitPathShort(vec.bufKE[:0], path, separator)
}

// A wrapper around bytealg.AppendSplitEntryString with additional logic for checking square brackets and "@" separator.
func (vec *Vector) appendSplitPathShort(dst []entry.Entry64, s, sep string) []entry.Entry64 {
	_, _ = splitTable[math.MaxUint8], splitDelta[math.MaxUint8]
	n, m := uint32(len(s)), len(sep)
	if n == 0 {
		return dst
	}

	_ = s[n-1]
	var lo, hi uint32
	for lo < n {
		for i := lo; i < n; i++ {
			switch {
			case splitTable[s[i]]:
				hi = i
				if e := entry.NewEntry64(lo, hi); lo != hi {
					dst = append(dst, e)
				}
				lo = i + splitDelta[s[i]]
			case m == 1 && s[i] == sep[0]:
				hi = i
				if e := entry.NewEntry64(lo, hi); lo != hi {
					dst = append(dst, e)
				}
				lo = i + 1
			case m > 1 && strings.HasPrefix(s[i:], sep):
				hi = i
				if e := entry.NewEntry64(lo, hi); lo != hi {
					dst = append(dst, e)
				}
				lo = i + uint32(len(sep))
			}
			if i == n-1 {
				hi = n
				if e := entry.NewEntry64(lo, hi); lo != hi && !splitTable[s[i]] {
					dst = append(dst, e)
				}
				goto exit
			}
		}
	}

exit:
	return dst
}

func (vec *Vector) appendSplitPath(dst []entry.Entry64, s, sep string) []entry.Entry64 {
	var t indextoken.Tokenizer[string]
	for {
		lo, hi := t.NextLH(s)
		if lo == hi {
			break
		}
		e := entry.NewEntry64(uint32(lo), uint32(hi))
		dst = append(dst, e)
	}
	return dst
}

var (
	splitTable = [math.MaxUint8 + 1]bool{}
	splitDelta = [math.MaxUint8 + 1]uint32{}
)

func init() {
	splitTable['@'] = true
	splitTable['['] = true
	splitTable[']'] = true

	splitDelta['['] = 1
	splitDelta[']'] = 1
}
