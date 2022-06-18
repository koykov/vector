# Vector API

Provide Vector API for various parsers.

The main idea: many data interchange formats (such as JSON, XML, ...) stores data as a tree.
Parsers of that formats reproduces that tree in a memory somehow or other. 

This parser uses different way: it stores all parsed nodes (key-value pairs) in a special array (vector).
That way protects from redundant memory allocations and reduces pointers.
In fact, vector contains only two pointers (array of nodes and array of indexes).
GC omits checking of that type of structs.

### Comparison

All known parsers has the following data structures (approximately) to represent node value (key-value pair):
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
As you see during parsing will be produced tons of pointers. Better consider this with an example of JSON:
```json
{
  "a": true,
  "b": {
    "c": "foo",
    "d": [
      5,
      3.1415,
      812.48927
    ]
  }
}
```

Majority of JSON parsers will build array of nodes like:

| 0               | 1          | 2               | 3           | 4                   | 5           | 6           | 7              |
|-----------------|------------|-----------------|-------------|---------------------|-------------|-------------|----------------|
| type: obj       | type: bool | type: obj       | type: str   | type: arr           | type: num   | type: num   | type: num      |
| key: ""         | key: "a"   | key: ""         | key: "c"    | key: ""             | key: ""     | key: ""     | key: ""        |
| str: ""         | str: ""    | str: ""         | str: "foo"  | str: ""             | str: ""     | str: ""     | str: ""        |
| num: 0          | num: 0     | num: 0          | num: 0      | num: 0              | num: 5      | num: 3.1415 | num: 812.48927 |
| bool: false     | bool: true | bool: false     | bool: false | bool: false         | bool: false | bool: false | bool: false    |
| child: [*1, *2] | child: []  | child: [*3, *4] | child: []   | child: [*5, *6, *7] | child: []   | child: []   | child: []      |

As you can see, independent of JSON node type, each parsed node contains at least 3 pointers:
* key (string)
* str (string)
* child (slice of node pointers)

JSON vector has different approach and build the following array of nodes and index:

Vector:

| 0          | 1          | 2           | 3           | 4           | 5           | 6           | 7           |
|------------|------------|-------------|-------------|-------------|-------------|-------------|-------------|
| type: obj  | type: bool | type: obj   | type: str   | type: arr   | type: num   | type: num   | type: num   |
| key pos: 0 | key pos: 5 | key pos: 18 | key pos: 29 | key pos: 45 | key pos: 0  | key pos: 0  | key pos: 0  |
| key len: 0 | key len: 1 | key len: 1  | key len: 1  | key len: 1  | key len: 0  | key len: 0  | key len: 0  |
| val pos: 0 | val pos: 9 | val pos: 0  | val pos: 34 | val pos: 0  | val pos: 57 | val pos: 64 | val pos: 80 |
| val len: 0 | val len: 4 | val len: 0  | val len: 3  | val len: 0  | val len: 1  | val len: 6  | val len: 9  |
| depth: 0   | depth: 1   | depth: 1    | depth: 2    | depth: 2    | depth: 3    | depth: 3    | depth: 3    |
| idx pos: 0 | idx pos: 0 | idx pos: 0  | idx pos: 0  | idx pos: 0  | idx pos: 0  | idx pos: 0  | idx pos: 0  |
| idx len: 2 | idx len: 0 | idx len: 2  | idx len: 0  | idx len: 3  | idx len: 0  | idx len: 0  | idx len: 0  |

Index (Y-axis means depth, X-axis means position in the index):

| X/Y | 0 | 1 | 2 |
|-----|---|---|---|
| 0   | 0 | - | - |
| 1   | 1 | 2 | - |
| 2   | 3 | 4 | - |
| 3   | 5 | 6 | 7 |

Each node contains only int variables and therefore avoid escapes to heap at all.
The whole vector contains only two pointers - nodes array and index matrix. GC checks it instantly.
That way also allows parsing JSON with unlimited depths.
