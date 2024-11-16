package vector

import "io"

type Interface interface {
	// SetHelper provides Helper to escape/unescape strings.
	SetHelper(helper Helper)

	// Parse parses source bytes.
	Parse(source []byte) error
	// ParseCopy makes a copy of source bytes and parse it.
	ParseCopy(source []byte) error
	// ParseString parses source string.
	// Note, string parsing may be unsafe due to in-place unescape issues. Use ParseCopyString or Parse in that case.
	ParseString(source string) error
	// ParseCopyString makes a copy of string and parse it.
	// This method is safe to parse immutable strings.
	ParseCopyString(source string) error
	// ParseStr is a legacy version of ParseString.
	// DEPRECATED: use ParseString instead.
	ParseStr(source string) error
	// ParseCopyStr is a legacy version of ParseCopyString.
	// DEPRECATED: use ParseCopyString instead.
	ParseCopyStr(source string) error

	// ParseReader takes source from r and parse it.
	ParseReader(r io.Reader) error

	// Root returns first root node.
	Root() *Node
	// RootByIndex returns root node by given index. If index overflows count of root nodes, the NULL node will return.
	RootByIndex(idx int) *Node
	// RootTop returns last root node.
	RootTop() *Node
	// Each applies closure to each root node.
	Each(fn func(idx int, fn *Node))
	// Exists checks if root node contains a child with given key.
	// The NULL node will return if node doesn't exist.
	Exists(key string) bool

	// Getters group.
	// Note, the NULL node will return if node doesn't exist by given keys.

	// Get returns node by given keys.
	Get(keys ...string) *Node
	// GetObject looks and get object node by given keys.
	GetObject(keys ...string) *Node
	// GetArray looks and get array node by given keys.
	GetArray(keys ...string) *Node
	// GetBytes looks and get bytes value by given keys.
	GetBytes(keys ...string) []byte
	// GetString looks and get string value by given keys.
	GetString(keys ...string) string
	// GetBool looks and get bool value by given keys.
	GetBool(keys ...string) bool
	// GetFloat looks and get float value by given keys.
	GetFloat(keys ...string) (float64, error)
	// GetInt looks and get integer value by given keys.
	GetInt(keys ...string) (int64, error)
	// GetUint looks and get unsigned integer value by given keys.
	GetUint(keys ...string) (uint64, error)

	// Getters by path/separator (PS) group.
	// Note, the NULL node will return if node doesn't exist by given keys.

	// GetPS returns node by given path and separator.
	GetPS(path, separator string) *Node
	// GetObjectPS looks and get object node by given path and separator.
	GetObjectPS(path, separator string) *Node
	// GetArrayPS looks and get array node by given path and separator.
	GetArrayPS(path, separator string) *Node
	// GetBytesPS looks and get bytes value by given path and separator.
	GetBytesPS(path, separator string) []byte
	// GetStringPS looks and get string value by given path and separator.
	GetStringPS(path, separator string) string
	// GetBoolPS looks and get bool value by given path and separator.
	GetBoolPS(path, separator string) bool
	// GetFloatPS looks and get float value by given path and separator.
	GetFloatPS(path, separator string) (float64, error)
	// GetIntPS looks and get integer value by given path and separator.
	GetIntPS(path, separator string) (int64, error)
	// GetUintPS looks and get unsigned integer value by given path and separator.
	GetUintPS(path, separator string) (uint64, error)

	// Dot-getters (the same as PS getters but with hardcoded dot (".") separator).

	// Dot looks and get node by given path and "." separator.
	Dot(path string) *Node
	// DotObject looks and get object node by given path and "." separator.
	DotObject(path string) *Node
	// DotArray looks and get array node by given path and "." separator.
	DotArray(path string) *Node
	// DotBytes looks and get bytes value by given path and "." separator.
	DotBytes(path string) []byte
	// DotString looks and get string value by given path and "." separator.
	DotString(path string) string
	// DotBool looks and get bool value by given path and "." separator.
	DotBool(path string) bool
	// DotFloat looks and get float value by given path and "." separator.
	DotFloat(path string) (float64, error)
	// DotInt looks and get integer value by given path and "." separator.
	DotInt(path string) (int64, error)
	// DotUint looks and get unsigned integer value by given path and "." separator.
	DotUint(path string) (uint64, error)

	// KeepPtr guarantees that vector object wouldn't be collected by GC.
	KeepPtr()

	// Beautify formats first root node in human-readable representation to w.
	Beautify(w io.Writer) error
	// Marshal serializes first root node to w.
	Marshal(w io.Writer) error

	// ErrorOffset returns last error offset.
	ErrorOffset() int
	// Prealloc prepares space for further parse.
	Prealloc(size uint)
	// Reset vector data.
	Reset()
}
