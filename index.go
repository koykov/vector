package vector

// Index represents nodes index matrix.
//
// Contain indexes of nodes in the vector divided by depth.
// Y-axis means depth, X-axis means position in index.
type Index struct {
	// Index tree.
	tree [][]int
	// Index depth.
	depth int
}

// Register new index for given depth.
func (idx *Index) Register(depth, i int) int {
	if len(idx.tree) <= depth {
		for len(idx.tree) <= depth {
			idx.tree = append(idx.tree, nil)
			idx.depth = len(idx.tree)
		}
	}
	idx.tree[depth] = append(idx.tree[depth], i)
	return len(idx.tree[depth])
}

// Len returns length of index row registered on depth.
func (idx *Index) Len(depth int) int {
	if len(idx.tree) <= depth {
		return 0
	}
	return len(idx.tree[depth])
}

// GetRow returns indices row registered at given depth.
func (idx *Index) GetRow(depth int) []int {
	if depth < 0 || depth >= len(idx.tree) {
		return nil
	}
	return idx.tree[depth]
}

// Reset rest of the index starting of given depth and offset in the tree.
func (idx *Index) Reset(depth, offset int) {
	if depth >= len(idx.tree) {
		return
	}
	if len(idx.tree[depth]) > offset {
		idx.tree[depth] = idx.tree[depth][:offset]
	}
	if depth+1 < len(idx.tree) {
		for i := depth + 1; i < len(idx.tree); i++ {
			idx.tree[i] = idx.tree[i][:0]
		}
	}
}

// Get subset [s:e] of index row registered on depth.
func (idx *Index) get(depth, s, e int) []int {
	l := idx.Len(depth)
	if l > s {
		return idx.tree[depth][s:e]
	}
	return nil
}

// Get index value.
func (idx *Index) val(depth, i int) int {
	return idx.tree[depth][i]
}

// Reset index object.
func (idx *Index) reset() {
	for i := 0; i < len(idx.tree); i++ {
		idx.tree[i] = idx.tree[i][:0]
	}
	idx.depth = 0
}
