package vector

// Index represents nodes index matrix.
//
// Contain indexes of nodes in the vector divided by depth.
// Y-axis means depth, X-axis means position in index.
type Index struct {
	// Index tree.
	tree []branch
	// Index depth.
	depth uint32
}

type branch struct {
	buf []uint32
	len uint32
}

// Register new index for given depth.
func (idx *Index) Register(depth, i uint32) uint32 {
	if idx.ulen() <= depth {
		for idx.ulen() <= depth {
			idx.tree = append(idx.tree, branch{})
			idx.depth = idx.ulen()
		}
	}
	b := &idx.tree[depth]
	if b.len < uint32(len(b.buf)) {
		b.buf[b.len] = i
	} else {
		b.buf = append(b.buf, i)
	}
	b.len++
	return b.len
}

// Len returns length of index row registered on depth.
func (idx *Index) Len(depth uint32) uint32 {
	if idx.ulen() <= depth {
		return 0
	}
	return idx.tree[depth].len
}

// GetRow returns indices row registered at given depth.
func (idx *Index) GetRow(depth uint32) []uint32 {
	if depth < 0 || depth >= idx.ulen() {
		return nil
	}
	b := &idx.tree[depth]
	return b.buf[:b.len]
}

// Reset rest of the index starting of given depth and offset in the tree.
func (idx *Index) Reset(depth, offset uint32) {
	if depth >= idx.ulen() {
		return
	}
	if idx.tree[depth].len > offset {
		idx.tree[depth].buf = idx.tree[depth].buf[:offset]
	}
	if depth+1 < idx.ulen() {
		for i := depth + 1; i < idx.ulen(); i++ {
			idx.tree[i].len = 0
		}
	}
}

// Get subset [s:e] of index row registered on depth.
func (idx *Index) get(depth, s, e uint32) []uint32 {
	l := idx.Len(depth)
	if l > s {
		return idx.tree[depth].buf[s:e]
	}
	return nil
}

// Get index value.
func (idx *Index) val(depth, i uint32) uint32 {
	return idx.tree[depth].buf[i]
}

// Reset index object.
func (idx *Index) reset() {
	for i := 0; i < len(idx.tree); i++ {
		idx.tree[i].len = 0
	}
	idx.depth = 0
}

func (idx *Index) ulen() uint32 {
	return uint32(len(idx.tree))
}
