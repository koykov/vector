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
