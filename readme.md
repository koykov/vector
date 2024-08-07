# Vector API

Provide Vector API for vector parsers.

The main idea: popular data interchange formats (such as JSON, XML, ...) stores data as a tree.
Parsers of that formats reproduces that tree in a memory somehow or other. This makes a lot of pointers in the heap and
GC does a lot of work during marking. Moreover, each of them makes a copy of source data, it's often is redundant. 

This parser uses different way: it stores all parsed nodes (key-value pairs) in a special array (vector). Instead of
pointers to child nodes, each node stores adjacency list of indices of childs. In fact (to reduce pointers) node stores
not an array of indices, but offset/length of child indices in special struct calls `Index` (see picture #2 below).

Thus, the main purpose of the whole project is minimizing of pointers and thereby cut the costs of GC work.
An additional purpose is a memory economy.

Let's check difference of that approaches on example. Source document:
```json
{
  "a":{"c":"foobar","d":3.1415},
  "b":["asdfgh","zxcvb"]
}
```

Typical parser will make a tree in memory like this:
<img width="100%" src="static/typical.svg" alt="">
> **_NOTE:_**  Each arrow represents a pointer, same as each non-empty string in nodes.    

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

Of course, this advantage has a cost - writing new parser using vector API is a hard challenge. Also, debug of vector
instance is non-trivial, due to debugger shows not a data (e.g strings), but offset/length of data in arrays outside.

## API

### Parsing

Vector API provides four methods to parsing:
```go
func (Vector) Parse([]byte) error
func (Vector) ParseCopy([]byte) error
func (Vector) ParseString(string) error
func (Vector) ParseCopyString(string) error
```

Copy-versions allow to make a copy of source data explicitly. By default, vector not makes a copy and nodes "points" to
memory outside of vector. It's a developer responsibility to extend a life of a source data at least the same as vector's
life. If it's impossible, the better way is to use a copy-methods.

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
They take a variadic path to field you require to read. Methods with exact types (`GetInt`, `GetBool`, ...) additionally
check a type of fields in source documents.

Vector API allows avoiding variadic usage using path with separator:
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

Vector API allows to do the opposite operation - compose original document from parsed data:
```go
func (Vector) Beautify(writer io.Writer) error
func (Vector) Marshal(writer io.Writer) error
```

`Beautify` method writer to writer a human-readable view of the document, `Marshal` - minimized version.

### Error handling

vector may return an error during parsing. The error may be impersonal, like "unexpected identifier" and provides no
information about the exact position in the document where error occurred. The following method may help in that case:
```go
func (Vector) ErrorOffset() int
```

### Iterating

If vector was used to parse more than one document, you may iterate them avoiding use of `RootByIndex` method:
```go
func (Vector) Each(fn func(index int, node *Node))
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
func (Node) Each(fn func(index int, node *Node))
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
func (Node) RemoveIf(cond func(index int, node *Node) bool)
```

### Child nodes access

```go
func (Node) Children() []Node
func (Node) ChildrenIndices() []int
```

### Serialization

Serialization is similar to vector API, but allows serializing only current node and its childrens (recursively):
```go
func (Node) Beautify(writer io.Writer) error
func (Node) Marshal(writer io.Writer) error
```

Thus, you may serialize not the whole object, but only necessary part of it.

## Helper

The important part of the API. It must realize the interface:
```go
type Helper interface {
	Indirect(ptr *Byteptr) []byte        // in-place unescape
	Beautify(writer io.Writer, node *Node) error
	Marshal(writer io.Writer, node *Node) error 
}
```

Since vector API is a common API for parsers, the concrete realization if that interface provides support for
concrete data formats ([JSON](https://github.com/koykov/jsonvector/blob/master/helper.go),
[XML](https://github.com/koykov/xmlvector/blob/master/helper.go), ...).

Note: unescaping (indirect) is happening in-place, not additional memory required.

## Pooling

Vector was implemented as high-load solution by design and therefore requires pooling to use: 
```go
vec := jsonvector.Acquire()
defer jsonvector.Release(vec)
// or
vec := xmlvector.Acquire()
defer xmlvector.Release(vec)
// ...
```

Vector package doesn't provide pooling (because it is a common API), but high level packages provides `Acquire`/`Release`
methods (working with `sync.Pool` implicitly).

In fact, you may use vectors not exactly in high-loaded project. In that case pooling may be omitted:
```go
vec := jsonvector.NewVector()
// or
vec := xmlvector.NewVector()
// ...
```

## Performance

There are versus-projects for each of the existing vector parsers:
* https://github.com/koykov/versus/tree/master/jsonvector
* https://github.com/koykov/versus/tree/master/xmlvector
* https://github.com/koykov/versus/tree/master/urlvector
* https://github.com/koykov/versus/tree/master/halvector
* todo: yamlvector

The most interesting is [jsonvector](https://github.com/koykov/versus/tree/master/jsonvector), due to worthy opponents.
Thus, let's focus on jsonvector.

### Testing dataset
* [small.json](https://github.com/koykov/versus/blob/master/jsonvector/testdata/small.json)
* [medium.json](https://github.com/koykov/versus/blob/master/jsonvector/testdata/medium.json)
* [large.json](https://github.com/koykov/versus/blob/master/jsonvector/testdata/large.json)
* [canada.json](https://github.com/koykov/versus/blob/master/jsonvector/testdata/canada.json)
* [citm_catalog.json](https://github.com/koykov/versus/blob/master/jsonvector/testdata/citm_catalog.json)
* [twitter.json](https://github.com/koykov/versus/blob/master/jsonvector/testdata/twitter.json)

### Testing stages
* [i7-7700HQ](https://github.com/koykov/versus/blob/master/jsonvector/benchstat/i7-7700HQ--10n-5m--new.txt)
* [i7-1185G7](https://github.com/koykov/versus/blob/master/jsonvector/benchstat/i7-1185G7--10n-1m--new.txt)
* [Apple M2](https://github.com/koykov/versus/blob/master/jsonvector/benchstat/Apple-M2--10n-1m--new.txt)
* [Xeon 4214](https://github.com/koykov/versus/blob/master/jsonvector/benchstat/Xeon-4214--10n-1m--new.txt)

All stages works on Ubuntu 22.04, Go version 1.22.

### Competitors
* https://github.com/valyala/fastjson
* https://github.com/koykov/jsonvector
* https://github.com/minio/simdjson-go

### The result

Server CPU Xeon 4214 is the most interested due to its closeness to production conditions:

#### sec/op:
| DS/lib       | fastjson         | jsonvector       | simdjson         |
|--------------|------------------|------------------|------------------|
| small.json   | 31.61n ± 10%     | **30.21n ±  7%** | 162.2n ±  7%     | 
| medium.json  | **214.0n ±  9%** | 221.7n ±  2%     | 590.9n ±  8%     |
| large.json   | **2.452µ ±  1%** | 3.168µ ±  2%     | 5.910µ ± 14%     |
| canada.json  | 1.087m ±  3%     | 1.253m ±  0%     | **948.1µ ±  1%** |
| citm.json    | 375.2µ ±  1%     | 280.6µ ±  1%     | **212.0µ ±  4%** |
| twitter.json | 101.88µ ±  3%    | **80.80µ ±  1%** | 114.8µ ± 10%     |

#### B/s
| DS/lib       | fastjson          | jsonvector        | simdjson          |
|--------------|-------------------|-------------------|-------------------|
| small.json   | 5.598Gi ± 11%     | **5.858Gi ±  8%** | 1.091Gi ±  6%     | 
| medium.json  | **10.14Gi ± 10%** | 9.784Gi ±  2%     | 3.671Gi ±  8%     |
| large.json   | **10.68Gi ±  1%** | 8.267Gi ±  2%     | 4.431Gi ± 13%     |
| canada.json  | 1.928Gi ±  3%     | 1713.4Mi ±  0%    | **2.211Gi ±  1%** |
| citm.json    | 4.287Gi ±  1%     | 5.732Gi ±  1%     | **7.587Gi ±  4%** |
| twitter.json | 5.773Gi ±  3%     | **7.279Gi ±  1%** | 5.124Gi ±  9%     |

Results are similar to each other, and this is a good achievement for vector, because the main advantage of vector -
long work outside of synthetic tests and optimizations for pointers. The similar speed is an extra bonus.
