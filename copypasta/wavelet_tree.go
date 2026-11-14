package copypasta

/*
Wavelet Tree (小波树)
https://codeforces.com/blog/entry/52854
https://oi-wiki.org/ds/wavelet-tree/

Static structure on integer array with values in [lo, hi].
Build: O(n log U).  Query: O(log U).  U = hi−lo+1.

Operations:
  Kth(l, r, k)     — k-th smallest (1-indexed k) in a[l..r]
  CountLE(l, r, v) — count of elements ≤ v in a[l..r]

Related problems:
https://www.luogu.com.cn/problem/P3834
https://codeforces.com/problemset/problem/786/C  2400
https://atcoder.jp/contests/abc234/tasks/abc234_d
*/

type waveletNode struct {
	lo, hi      int
	b           []int // b[i] = # elements in a[0..i-1] that went to left child
	left, right *waveletNode
}

// newWaveletTree builds a wavelet tree for a (values in [lo, hi]).
// The slice a is rearranged in place; pass a copy to preserve the original.
func newWaveletTree(a []int, lo, hi int) *waveletNode {
	n := len(a)
	t := &waveletNode{lo: lo, hi: hi, b: make([]int, n+1)}
	if lo == hi || n == 0 {
		return t
	}
	mid := lo + (hi-lo)/2
	tmp := make([]int, n)
	j := 0
	for i, v := range a {
		if v <= mid {
			t.b[i+1] = t.b[i] + 1
			tmp[j] = v
			j++
		} else {
			t.b[i+1] = t.b[i]
		}
	}
	leftCnt := j
	for _, v := range a {
		if v > mid {
			tmp[j] = v
			j++
		}
	}
	copy(a, tmp)
	t.left = newWaveletTree(a[:leftCnt], lo, mid)
	t.right = newWaveletTree(a[leftCnt:], mid+1, hi)
	return t
}

// Kth returns the k-th smallest element (1-indexed) in a[l..r].
func (t *waveletNode) Kth(l, r, k int) int {
	if t.lo == t.hi {
		return t.lo
	}
	lb := t.b[r+1] - t.b[l] // elements in a[l..r] that went to left child
	if k <= lb {
		return t.left.Kth(t.b[l], t.b[r+1]-1, k)
	}
	return t.right.Kth(l-t.b[l], r-t.b[r+1], k-lb)
}

// CountLE returns the number of elements ≤ v in a[l..r].
func (t *waveletNode) CountLE(l, r, v int) int {
	if l > r || v < t.lo {
		return 0
	}
	if v >= t.hi {
		return r - l + 1
	}
	if t.lo == t.hi {
		return r - l + 1
	}
	mid := t.lo + (t.hi-t.lo)/2
	lb := t.b[r+1] - t.b[l]
	if v <= mid {
		return t.left.CountLE(t.b[l], t.b[r+1]-1, v)
	}
	return lb + t.right.CountLE(l-t.b[l], r-t.b[r+1], v)
}
