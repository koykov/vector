package vector

import "io"

type Interface interface {
	SetHelper(helper Helper)

	Parse([]byte) error
	ParseCopy([]byte) error
	ParseString(string) error
	ParseCopyString(string) error
	// DEPRECATED: use ParseString instead.
	ParseStr(string) error
	// DEPRECATED: use ParseCopyString instead.
	ParseCopyStr(string) error

	Root() *Node
	RootByIndex(int) *Node
	RootTop() *Node
	Each(fn func(int, *Node))
	Exists(string) bool

	Get(...string) *Node
	GetObject(...string) *Node
	GetArray(...string) *Node
	GetBytes(...string) []byte
	GetString(...string) string
	GetBool(...string) bool
	GetFloat(...string) (float64, error)
	GetInt(...string) (int64, error)
	GetUint(...string) (uint64, error)

	GetPS(string, string) *Node
	GetObjectPS(string, string) *Node
	GetArrayPS(string, string) *Node
	GetBytesPS(string, string) []byte
	GetStringPS(string, string) string
	GetBoolPS(string, string) bool
	GetFloatPS(string, string) (float64, error)
	GetIntPS(string, string) (int64, error)
	GetUintPS(string, string) (uint64, error)

	Dot(string) *Node
	DotObject(string) *Node
	DotArray(string) *Node
	DotBytes(string) []byte
	DotString(string) string
	DotBool(string) bool
	DotFloat(string) (float64, error)
	DotInt(string) (int64, error)
	DotUint(string) (uint64, error)

	KeepPtr()

	Beautify(io.Writer) error
	Marshal(io.Writer) error

	ErrorOffset() int
	Prealloc(uint)
	Reset()
}
