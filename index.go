package vector

type index struct {
	tree  [][]int
	depth int
}

func (idx *index) register(depth, i int) int {
	if len(idx.tree) <= depth {
		for len(idx.tree) <= depth {
			idx.tree = append(idx.tree, nil)
			idx.depth = len(idx.tree)
		}
	}
	idx.tree[depth] = append(idx.tree[depth], i)
	return len(idx.tree[depth])
}

func (idx *index) len(depth int) int {
	if len(idx.tree) <= depth {
		return 0
	}
	return len(idx.tree[depth])
}

func (idx *index) get(depth, s, e int) []int {
	l := idx.len(depth)
	if l > s {
		return idx.tree[depth][s:e]
	}
	return nil
}

func (idx *index) val(depth, i int) int {
	return idx.tree[depth][i]
}

func (idx *index) reset() {
	for i := 0; i < len(idx.tree); i++ {
		idx.tree[i] = idx.tree[i][:0]
	}
	idx.depth = 0
}
