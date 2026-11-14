package copypasta

/*
Skew Heap (斜堆)
https://en.wikipedia.org/wiki/Skew_heap
https://oi-wiki.org/ds/leftist-tree/

A self-adjusting mergeable min-heap. Simpler than a leftist tree:
every merge swaps left/right children unconditionally.

Operations (amortized):
  skewMerge(a, b)  — merge two heaps → O(log n)
  skewPush(h, v)   — insert value   → O(log n)
  skewPop(h)       — remove min     → O(log n)

Use cases: Dijkstra with lazy deletions, Prim on dense graphs,
           offline dynamic MST, Huffman tree construction.

Related problems:
https://www.luogu.com.cn/problem/P3377  (leftist tree template)
https://codeforces.com/problemset/problem/1439/C  2100
*/

type skewNode struct {
	val         int
	left, right *skewNode
}

// skewMerge merges two min-heaps. Either may be nil.
func skewMerge(a, b *skewNode) *skewNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a.val > b.val {
		a, b = b, a
	}
	a.right = skewMerge(a.right, b)
	a.left, a.right = a.right, a.left
	return a
}

// skewPush inserts val into heap h and returns the new root.
func skewPush(h *skewNode, val int) *skewNode {
	return skewMerge(h, &skewNode{val: val})
}

// skewPop removes the minimum from heap h and returns the new root.
// h must be non-nil.
func skewPop(h *skewNode) *skewNode {
	return skewMerge(h.left, h.right)
}
