package vector

import "github.com/koykov/bytealg"

type Flag int

const (
	FlagEscape Flag = iota
)

type Byteptr struct {
	bytealg.Byteptr
	flagEsc bool
}

func (m *Byteptr) RawBytes() []byte {
	return m.Bytes()
}

func (m *Byteptr) SetFlag(flag Flag, value bool) {
	switch flag {
	case FlagEscape:
		m.flagEsc = value
	}
}

func (m *Byteptr) GetFlag(flag Flag) bool {
	switch flag {
	case FlagEscape:
		return m.flagEsc
	}
	return false
}

func (m *Byteptr) Reset() {
	m.Byteptr.Reset()
	m.flagEsc = false
}
