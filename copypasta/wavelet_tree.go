package copypasta

import "sort"

/*
Wavelet Tree (小波树)
https://codeforces.com/blog/entry/52854
https://oi-wiki.org/ds/wavelet-tree/

Static structure on an integer slice supporting range order-statistic queries.
Automatically coordinate-compresses distinct values so arbitrary int values
(e.g. [-1e18, 1e18]) take O(n log sigma) build time/space and O(log sigma) per query,
where sigma <= n is the number of distinct values.

Operations:
  Kth(l, r, k)            — k-th smallest (1-indexed k) in a[l..r]
  CountLE(l, r, v)        — count of elements <= v in a[l..r]
  CountRange(l, r, lo, hi)— count of elements in [lo, hi] in a[l..r]

Related problems:
https://www.luogu.com.cn/problem/P3834
https://codeforces.com/problemset/problem/786/C  2400
https://atcoder.jp/contests/abc234/tasks/abc234_d
*/

type waveletNode struct {
	lo, hi      int
	b           []int // b[i] = # elements in current slice[0..i-1] routed to left child
	left, right *waveletNode
}

type waveletTree struct {
	root *waveletNode
	vals []int // sorted unique values
}

func buildWaveletRanks(ranks, tmp []int, lo, hi int) *waveletNode {
	n := len(ranks)
	t := &waveletNode{lo: lo, hi: hi}
	if lo == hi || n == 0 {
		return t
	}
	t.b = make([]int, n+1)
	mid := lo + int(uint(hi-lo)>>1)
	j := 0
	for i, v := range ranks {
		if v <= mid {
			t.b[i+1] = t.b[i] + 1
			tmp[j] = v
			j++
		} else {
			t.b[i+1] = t.b[i]
		}
	}
	leftCnt := j
	for _, v := range ranks {
		if v > mid {
			tmp[j] = v
			j++
		}
	}
	copy(ranks, tmp[:n])
	t.left = buildWaveletRanks(ranks[:leftCnt], tmp[:leftCnt], lo, mid)
	t.right = buildWaveletRanks(ranks[leftCnt:], tmp[leftCnt:], mid+1, hi)
	return t
}

// newWaveletTree builds a coordinate-compressed wavelet tree over a without mutating a.
func newWaveletTree(a []int) *waveletTree {
	n := len(a)
	if n == 0 {
		return &waveletTree{}
	}
	sorted := append([]int(nil), a...)
	sort.Ints(sorted)
	vals := sorted[:0]
	for _, v := range sorted {
		if len(vals) == 0 || vals[len(vals)-1] != v {
			vals = append(vals, v)
		}
	}
	ranks := make([]int, n)
	for i, v := range a {
		ranks[i] = sort.SearchInts(vals, v)
	}
	tmp := make([]int, n)
	root := buildWaveletRanks(ranks, tmp, 0, len(vals)-1)
	return &waveletTree{root: root, vals: vals}
}

func (t *waveletNode) kthRank(l, r, k int) int {
	if t.lo == t.hi {
		return t.lo
	}
	lb := t.b[r+1] - t.b[l]
	if k <= lb {
		return t.left.kthRank(t.b[l], t.b[r+1]-1, k)
	}
	return t.right.kthRank(l-t.b[l], r-t.b[r+1], k-lb)
}

func (t *waveletNode) countLERank(l, r, v int) int {
	if l > r || v < t.lo {
		return 0
	}
	if v >= t.hi || t.lo == t.hi {
		return r - l + 1
	}
	mid := t.lo + int(uint(t.hi-t.lo)>>1)
	lb := t.b[r+1] - t.b[l]
	if v <= mid {
		return t.left.countLERank(t.b[l], t.b[r+1]-1, v)
	}
	return lb + t.right.countLERank(l-t.b[l], r-t.b[r+1], v)
}

// Kth returns the k-th smallest element (1-indexed) in a[l..r].
func (wt *waveletTree) Kth(l, r, k int) int {
	if l > r || k < 1 || k > r-l+1 || wt.root == nil {
		panic("waveletTree.Kth: invalid range or k")
	}
	return wt.vals[wt.root.kthRank(l, r, k)]
}

// CountLE returns the number of elements <= v in a[l..r].
func (wt *waveletTree) CountLE(l, r, v int) int {
	if l > r || len(wt.vals) == 0 {
		return 0
	}
	idx := sort.Search(len(wt.vals), func(i int) bool { return wt.vals[i] > v }) - 1
	if idx < 0 {
		return 0
	}
	if idx == len(wt.vals)-1 {
		return r - l + 1
	}
	return wt.root.countLERank(l, r, idx)
}

// CountRange returns the number of elements in [lo, hi] within a[l..r].
func (wt *waveletTree) CountRange(l, r, lo, hi int) int {
	if l > r || lo > hi || len(wt.vals) == 0 {
		return 0
	}
	ans := wt.CountLE(l, r, hi)
	if lo > wt.vals[0] {
		ans -= wt.CountLE(l, r, lo-1)
	}
	return ans
}
