# Vector API

Provide Vector API for vector parsers.

The main idea: popular data interchange formats (such as JSON, XML, ...) stores data as a tree.
Parsers of that formats reproduces that tree in a memory somehow or other. Moreover, each of them makes a copy of source
data, it's often is redundant. This makes a lot of pointers in the heap and GC does a lot of work during marking.

This parser uses different way: it stores all parsed nodes (key-value pairs) in a special array (vector). Instead of
pointers to child nodes, each node stores adjacency list of indices of childs. In fact, to reduce pointers node stores
not and array of indices, but offset/length of child indices in special struct calls `Index` (see below).

Thus, the main purpose of the whole project is minimising of pointers and thereby cut the costs of GC work.
An additional purpose if memory economy.

Let's check difference of that approaches on example:
Source document:
```json
{
  "a":{"c":"foobar","d":3.1415},
  "b":["asdfgh","zxcvb"]
}
```

Typical parse will make a tree in memory like this:
<img width="100%" src="static/typical.svg" alt="">
> **_NOTE:_**  Each arrow represents a pointer, same as an each non-empty string in nodes.    

, each node will an instance of struct like this:
```go
type Node struct {
	typ Type    // [null, object, array, string, number, true, false]
	obj Object  // one pointer to slice and N*2 pointers inside KeyValue struct, see below
	arr []*Node // one pointer for the slice and N pointers for each array item
	str string  // one pointer
}

type Object []KeyValue

type KeyValue struct {
	key string // one pointer
	val *Node  // one pointer to node
}
```

That way produces a too many pointers and on complex source document that count will grow even more.

Vector in memory will look like:
<img width="100%" src="static/vector.svg" alt="">
> **_NOTE:_**  In opposite to previous picture, arrows represents not a pointers, but integer indices.

Looks redundant and strange, isn't it? But that way allow for any type of source document produces constant number of
pointers:
* one for array of nodes
* one for index
* one for each row in index

## API

### Parsing

Vector API provides four methods to parsing:
```go
func (Vector) Parse([]byte) error
func (Vector) ParseCopy([]byte) error
func (Vector) ParseString(string) error
func (Vector) ParseCopyString(string) error
```

copy-versions allow to make a copy of source data explicitly. By default, vector not makes a copy and nodes "points" to
memory outside of vector. It's a developer responsibility to extend a life of a source data at least the same as vector's
life. If it's impossible, he better way is to use a copy-methods.

The exclusive feature of vector is a possibility to parse many source documents using one vector instance:
```go
vec.ParseString(`{"a":{"b":{"c":"foobar"}}}`)
vec.ParseString(`{"x":{"y":{"z":"asdfgh"}}}`)
s0 := vec.RootByIndex(0).DotString("a.b.c")
s1 := vec.RootByIndex(1).DotString("x.y.z")
println(s0, s1) // foobar asdfgh
```

Thus, vector minimizes pointers count on multiple source data as if parse only one source document.

### Reading

The basic reading methods:
```go
func (Vector) Get(path ...string) *Node
func (Vector) GetObject(path ...string) *Node
func (Vector) GetArray(path ...string) *Node
func (Vector) GetBytes(path ...string) []byte
func (Vector) GetString(path ...string) string
func (Vector) GetBool(path ...string) bool
func (Vector) GetFloat(path ...string) (float64, error)
func (Vector) GetInt(path ...string) (int64, error)
func (Vector) GetUint(path ...string) (uint64, error)
```
They takes a variadic path to field you require to read. Methods with exact types (`GetInt`, `GetBool`, ...) additionally
check a type of fields in source documents.

vector API allows to avoid using variadic using path with separator:
```go
func (Vector) GetPS(path, separator string) *Node
func (Vector) GetObjectPS(path, separator string) *Node
func (Vector) GetArrayPS(path, separator string) *Node
func (Vector) GetBytesPS(path, separator string) []byte
func (Vector) GetStringPS(path, separator string) string
func (Vector) GetBoolPS(path, separator string) bool
func (Vector) GetFloatPS(path, separator string) (float64, error)
func (Vector) GetIntPS(path, separator string) (int64, error)
func (Vector) GetUintPS(path, separator string) (uint64, error)
```
Example:
```go
vec.ParseString(`{"a":{"b":{"c":"foobar"}}}`)
s := vec.GetStringPS("a.b.c", ".")
println(s) // foobar
```

Due to most popular separator is a dot (".") there are special alias-methods:
```go
func (Vector) Dot(path string) *Node
func (Vector) DotObject(path string) *Node
func (Vector) DotArray(path string) *Node
func (Vector) DotBytes(path string) []byte
func (Vector) DotString(path string) string
func (Vector) DotBool(path string) bool
func (Vector) DotFloat(path string) (float64, error)
func (Vector) DotInt(path string) (int64, error)
func (Vector) DotUint(path string) (uint64, error)
```
Example:
```go
vec.ParseString(`{"a":{"b":{"c":"foobar"}}}`)
s := vec.DotString("a.b.c")
println(s) // foobar
```

### Serialization

vector API allows to do the opposite operation - compose original document from parsed data:
```go
func (Vector) Beautify(io.Writer) error
func (Vector) Marshal(io.Writer) error
```

`Beautify` method makes a human-readable view of the document, `Marshal` - minimized version.

### Error handling

vector may return an error during parsing. The error may be impersonal, like "unexpected identifier" and provides no
information about the exact position in the document where error occurred. The following method may help in that case:
```go
func (Vector) ErrorOffset() int
```

### Iterating

If vector was used to parse more than one document, you may iterate them avoiding use of `RootByIndex` method:
```go
func (Vector) Each(fn func(int, *Node))
```
Example:
```go
vec.ParseString(`{"a":{"b":{"c":"foobar"}}}`)
vec.ParseString(`{"x":{"y":{"z":"asdfgh"}}}`)
vec.Each(func(i int, node *Node) {
	node.Get("...")
})
```

## node API

### Reading

Similar to vector API, there are three groups of methods in node API:
* Get-methods
* GetPS-methods
* Dot-methods

In addition, node can return key/value separately as [Byteptr](byteptr.go) object:
```go
func (Node) Key() *Byteptr
func (Node) Value() *Byteptr
```
or directly convert them to exact types:
```go
// key
func (Node) KeyBytes() []byte
func (Node) KeyString() string
// value
func (Node) Bytes() []byte
func (Node) ForceBytes() []byte
func (Node) RawBytes() []byte
func (Node) String() string
func (Node) ForceString()
func (Node) Bool() bool
func (Node) Float() (float64, error)
func (Node) Int() (int64, error)
func (Node) Uint() (uint64, error)
func (Node) Object() *Node
func (Node) Array() *Node

func (Node) Type() Type
func (Node) Exists(key string) bool
```

### Iterating

If node has a type array or object, you may iterate through children nodes:
```go
func (Node) Each(fn func(idx int, node *Node))
```

### Sorting

Nodes of type array or object may be sorted by keys or values:
```go
func (Node) SortKeys() *Node // by keys
func (Node) Sort() *Node     // by values
```

### Removing

Node API supports predicating deletion:
```go
func (Node) RemoveIf(cond func(idx int, node *Node) bool)
```

### Child nodes access

```go
func (Node) Children() []Node
func (Node) ChildrenIndices() []int
```

### Serialization

Serialization is similar to vector API, but allows to serialize only current node and its childrens (recursively):
```go
func (Node) Beautify(io.Writer) error
func (Node) Marshal(io.Writer) error
```

Thus, you may serialize not the whole object, but only necessary part of it.

## Helper

The important part of the API. It must realize the interface:
```go
type Helper interface {
	Indirect(*Byteptr) []byte        // in-place unescape
	Beautify(io.Writer, *Node) error
	Marshal(io.Writer, *Node) error 
}
```

Since vector API is a common API for parsers, the concrete realization if that interface provides support for
concrete data formats ([JSON](https://github.com/koykov/jsonvector/blob/master/helper.go),
[XML](https://github.com/koykov/xmlvector/blob/master/helper.go), ...).

Note: unescaping (indirect) is happening in-place, not additional memory required.

## Pooling

vector was implemented as highload solution by design and therefore requires pooling to use: 
```go
vec := jsonvector.Acquire()
defer jsonvector.Release(vec)
// or
vec := xmlvector.Acquire()
defer xmlvector.Release(vec)
// ...
```

vector package doesn't provide pooling (because it is a common API), but high level packages provides `Acquire`/`Release`
methods (working with `sync.Pool` implicitly).

In fact, you may use vectors not exactly in high-loaded project. In that case pooling may be omitted:
```go
vec := jsonvector.NewVector()
// or
vec := xmlvector.NewVector()
// ...
```
