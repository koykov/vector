package vector

import "sync"

var testPool = sync.Pool{New: func() any { return &Vector{} }}
