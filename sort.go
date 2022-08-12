package vector

// Custom implementation of quick sort algorithm, special for type []Node.
// Need to avoid redundant allocation when using sort.Interface.
//
// sort.Interface problem:
// <code>
// type nodes []Node // type that implements sort.Interface
// ...
// children := node.Children() // get a slice of nodes to sort
// nodes := nodes(children)    // <- simple typecast, but produces an alloc (copy) due to taking address in the next line
// sort.Sort(&nodes)           // taking address
// ...
// </code>

const (
	sortByKey sortMode = iota
	sortByValue
)

type sortMode uint8

func pivot(p []Node, lo, hi int, mode sortMode) int {
	if len(p) == 0 {
		return 0
	}
	pi := &p[hi]
	i := lo - 1
	_ = p[len(p)-1]
	for j := lo; j <= hi-1; j++ {
		var mustSwap bool
		switch mode {
		case sortByKey:
			mustSwap = p[j].KeyString() < pi.KeyString()
		case sortByValue:
			mustSwap = p[j].String() < pi.String()
		}
		if mustSwap {
			i++
			p[i].SwapWith(&p[j])
		}
	}
	if i < hi {
		p[i+1].SwapWith(&p[hi])
	}
	return i + 1
}

func quickSort(p []Node, lo, hi int, mode sortMode) {
	if lo < hi {
		pi := pivot(p, lo, hi, mode)
		quickSort(p, lo, pi-1, mode)
		quickSort(p, pi+1, hi, mode)
	}
}
