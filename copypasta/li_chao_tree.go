package copypasta

import "math"

/*
Li Chao Tree (李超线段树)
https://cp-algorithms.com/geometry/li_chao-tree.html
https://oi-wiki.org/ds/li-chao-tree/

Maintains lines/segments y = kx+b and answers extremum queries at integer x in [lo, hi].

Time:
  - AddLine / QueryMax: O(log(hi-lo))
  - AddSegment:         O(log^2(hi-lo))
  - Merge (tree DP):    O(n log(hi-lo)) amortized across n lines

Space: O(q * log(hi-lo)) dynamic nodes

Related problems:
https://www.luogu.com.cn/problem/P4097  (segment Li Chao template)
LC1235 https://leetcode.cn/problems/maximum-profit-in-job-scheduling/
https://codeforces.com/problemset/problem/932/F   2400  (tree DP + Li Chao merge, min query)
https://atcoder.jp/contests/dp/tasks/dp_z
*/

type lcLine struct{ k, b int }

func (l lcLine) eval(x int) int { return l.k*x + l.b }

type lcNode struct {
	line        lcLine
	hasLine     bool
	left, right *lcNode
}

func lcMid(lo, hi int) int {
	return lo + int(uint(hi-lo)>>1)
}

// liChaoTree maintains lines/segments y=kx+b over x in [lo,hi] for max-y queries.
type liChaoTree struct {
	root   *lcNode
	lo, hi int
}

func newLiChaoTree(lo, hi int) *liChaoTree {
	return &liChaoTree{lo: lo, hi: hi}
}

func (t *liChaoTree) insert(node *lcNode, lo, hi int, line lcLine) *lcNode {
	if node == nil {
		return &lcNode{line: line, hasLine: true}
	}
	if !node.hasLine {
		node.line, node.hasLine = line, true
		return node
	}
	mid := lcMid(lo, hi)
	leftBetter := line.eval(lo) > node.line.eval(lo)
	midBetter := line.eval(mid) > node.line.eval(mid)
	if midBetter {
		node.line, line = line, node.line
	}
	if lo < hi {
		if leftBetter != midBetter {
			node.left = t.insert(node.left, lo, mid, line)
		} else {
			node.right = t.insert(node.right, mid+1, hi, line)
		}
	}
	return node
}

func (t *liChaoTree) insertSegment(node *lcNode, lo, hi, ql, qr int, line lcLine) *lcNode {
	if qr < lo || hi < ql || ql > qr {
		return node
	}
	if ql <= lo && hi <= qr {
		return t.insert(node, lo, hi, line)
	}
	if node == nil {
		node = &lcNode{}
	}
	mid := lcMid(lo, hi)
	if ql <= mid {
		node.left = t.insertSegment(node.left, lo, mid, ql, qr, line)
	}
	if qr > mid {
		node.right = t.insertSegment(node.right, mid+1, hi, ql, qr, line)
	}
	return node
}

// AddLine inserts the global line y = k*x + b over [lo, hi].
func (t *liChaoTree) AddLine(k, b int) {
	t.root = t.insert(t.root, t.lo, t.hi, lcLine{k, b})
}

// AddSegment inserts the segment y = k*x + b active only on integer x in [l, r].
func (t *liChaoTree) AddSegment(l, r, k, b int) {
	t.root = t.insertSegment(t.root, t.lo, t.hi, l, r, lcLine{k, b})
}

func (t *liChaoTree) mergeNodes(a, b *lcNode, lo, hi int) *lcNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if b.hasLine {
		a = t.insert(a, lo, hi, b.line)
	}
	if lo < hi {
		mid := lcMid(lo, hi)
		a.left = t.mergeNodes(a.left, b.left, lo, mid)
		a.right = t.mergeNodes(a.right, b.right, mid+1, hi)
	}
	return a
}

// Merge absorbs all lines from other into t (both must share [lo, hi]).
// other should not be used after Merge.
func (t *liChaoTree) Merge(other *liChaoTree) {
	if other == nil {
		return
	}
	t.root = t.mergeNodes(t.root, other.root, t.lo, t.hi)
}

func (t *liChaoTree) query(node *lcNode, lo, hi, x int) int {
	if node == nil {
		return math.MinInt
	}
	res := math.MinInt
	if node.hasLine {
		res = node.line.eval(x)
	}
	if lo == hi {
		return res
	}
	mid := lcMid(lo, hi)
	var child int
	if x <= mid {
		child = t.query(node.left, lo, mid, x)
	} else {
		child = t.query(node.right, mid+1, hi, x)
	}
	if child > res {
		return child
	}
	return res
}

// QueryMax returns the maximum y among all active lines/segments at x.
// Returns math.MinInt if no line covers x.
func (t *liChaoTree) QueryMax(x int) int {
	return t.query(t.root, t.lo, t.hi, x)
}

// liChaoMinTree wraps liChaoTree for min-y queries (e.g. CF 932F, AtCoder dp_z).
type liChaoMinTree struct {
	inner *liChaoTree
}

func newLiChaoMinTree(lo, hi int) *liChaoMinTree {
	return &liChaoMinTree{inner: newLiChaoTree(lo, hi)}
}

func (t *liChaoMinTree) AddLine(k, b int)          { t.inner.AddLine(-k, -b) }
func (t *liChaoMinTree) AddSegment(l, r, k, b int) { t.inner.AddSegment(l, r, -k, -b) }
func (t *liChaoMinTree) Merge(other *liChaoMinTree) {
	if other != nil {
		t.inner.Merge(other.inner)
	}
}

// QueryMin returns the minimum y at x, or math.MaxInt if no line covers x.
func (t *liChaoMinTree) QueryMin(x int) int {
	res := t.inner.QueryMax(x)
	if res == math.MinInt {
		return math.MaxInt
	}
	return -res
}
