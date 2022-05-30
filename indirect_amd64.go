package vector

// Indirect vector object from raw pointer.
//
// Fix "checkptr: pointer arithmetic result points to invalid allocation" error in race mode.
//go:noescape
func indirectVector1(_ uintptr) *Vector

//go:noescape
func indirectNode1(_ uintptr) *Node
