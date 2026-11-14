package copypasta

import "math"

/*
Li Chao Tree (李超线段树)
https://cp-algorithms.com/geometry/li_chao-tree.html
https://oi-wiki.org/ds/li-chao-tree/

Maintains a set of lines y = kx+b and answers max-y queries at integer x ∈ [lo, hi].

Time:  O(log(hi-lo)) per AddLine and QueryMax
Space: O(q·log(hi-lo)) dynamic nodes, q = number of inserted lines

Common application: Convex Hull Trick DP when x-queries are not monotone.

Related problems:
LC1235 https://leetcode.cn/problems/maximum-profit-in-job-scheduling/
https://codeforces.com/problemset/problem/932/F   2400
https://atcoder.jp/contests/dp/tasks/dp_z
*/

type lcLine struct{ k, b int }

func (l lcLine) eval(x int) int { return l.k*x + l.b }

type lcNode struct {
	line        lcLine
	hasLine     bool
	left, right *lcNode
}

// liChaoTree maintains lines y=kx+b over x∈[lo,hi] for max-y queries.
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
	mid := lo + (hi-lo)/2
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

// AddLine inserts the line y = k*x + b.
func (t *liChaoTree) AddLine(k, b int) {
	t.root = t.insert(t.root, t.lo, t.hi, lcLine{k, b})
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
	mid := lo + (hi-lo)/2
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

// QueryMax returns the maximum y among all inserted lines at x.
// Returns math.MinInt if no lines have been added.
func (t *liChaoTree) QueryMax(x int) int {
	return t.query(t.root, t.lo, t.hi, x)
}
