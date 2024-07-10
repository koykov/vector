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
<img width="100%" src="static/vector-color.svg" alt="">

Looks redundant and strange, isn't it? But that way allow for any type of source document produces constant number of
pointers:
* one for array of nodes
* one for index
* one for each row in index

In fact, node struct if more complex that present on the picture. Key/value isn't a bit-shifted lo/hi indices in source
data, but special struct [Byteptr](byteptr.go). It stores additionally raw (uintptr) "pointer" to source data and extra
flags. But the idea is the same: in GC eyes, both byteptr and node is simple structs and GC will not spend time to check
them.

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
