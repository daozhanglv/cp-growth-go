package copypasta

/*
Convex Hull Trick — Monotone (单调凸包优化 / 斜率优化DP)
https://cp-algorithms.com/geometry/convex_hull_trick.html
https://oi-wiki.org/dp/opt/slope/

Compared with Li Chao Tree (also in this package):
  - Li Chao Tree: arbitrary insert order, arbitrary query order, O(log n) each
  - CHT (this file): lines must be added in monotone slope order
      * chtDeque:  queries must also be monotone → O(1) amortized per query
      * chtStack:  arbitrary query order → O(log n) per query (binary search)

Common DP pattern:
  dp[i] = min over j<i of { a[j]*b[i] + c[j] }
  where a[j] is the "slope" and c[j] is the "intercept" of line j.

Related problems:
https://codeforces.com/problemset/problem/311/E  2500
https://codeforces.com/problemset/problem/319/C  2400
LC1235 https://leetcode.cn/problems/maximum-profit-in-job-scheduling/
https://atcoder.jp/contests/dp/tasks/dp_z
*/

// bad returns true if line b is never the maximum among a, b, c.
// i.e., the intersection of a,c is to the left of the intersection of a,b.
func chtBad(a, b, c chtLine) bool {
	// (c.b - a.b) * (a.k - b.k) >= (b.b - a.b) * (a.k - c.k)
	return (c.b-a.b)*(a.k-b.k) >= (b.b-a.b)*(a.k-c.k)
}

type chtLine struct{ k, b int }

func (l chtLine) eval(x int) int { return l.k*x + l.b }

// chtDeque supports:
//
//	AddLine: lines added in DECREASING k order
//	QueryMax: x must be INCREASING across queries
//
// Both operations are O(1) amortized.
//
// Lines are stored with index 0 = highest slope; new lines (lower slope) are
// appended at the back.  For increasing x queries, the optimal line moves from
// the back (low slope, optimal at small x) toward the front (high slope,
// optimal at large x), so we track a hi pointer that retreats from the back.
type chtDeque struct {
	lines []chtLine
	hi    int // back pointer: current candidate is lines[hi]; starts at len-1
}

func (c *chtDeque) AddLine(k, b int) {
	l := chtLine{k, b}
	for len(c.lines) >= 2 && chtBad(c.lines[len(c.lines)-2], c.lines[len(c.lines)-1], l) {
		c.lines = c.lines[:len(c.lines)-1]
	}
	c.lines = append(c.lines, l)
	c.hi = len(c.lines) - 1
}

// QueryMax returns max y at x. x must be >= all previous query x values.
func (c *chtDeque) QueryMax(x int) int {
	if len(c.lines) == 0 {
		panic("chtDeque: no lines added")
	}
	// Advance hi toward 0: if the line before hi is better at x, move hi back.
	for c.hi > 0 && c.lines[c.hi-1].eval(x) >= c.lines[c.hi].eval(x) {
		c.hi--
	}
	return c.lines[c.hi].eval(x)
}

// chtStack supports:
//
//	AddLine: lines added in DECREASING k order
//	QueryMax: arbitrary x (binary search)
type chtStack struct {
	lines []chtLine
}

func (c *chtStack) AddLine(k, b int) {
	l := chtLine{k, b}
	for len(c.lines) >= 2 && chtBad(c.lines[len(c.lines)-2], c.lines[len(c.lines)-1], l) {
		c.lines = c.lines[:len(c.lines)-1]
	}
	c.lines = append(c.lines, l)
}

// QueryMax returns max y at x (arbitrary order).
func (c *chtStack) QueryMax(x int) int {
	n := len(c.lines)
	if n == 0 {
		panic("chtStack: no lines added")
	}
	// binary search: find last line where lines[i] >= lines[i+1] at x
	lo, hi := 0, n-1
	for lo < hi {
		mid := (lo + hi) / 2
		if c.lines[mid].eval(x) >= c.lines[mid+1].eval(x) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return c.lines[lo].eval(x)
}
