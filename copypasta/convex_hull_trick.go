package copypasta

import "math/bits"

/*
Convex Hull Trick — Monotone (单调凸包优化 / 斜率优化DP)
https://cp-algorithms.com/geometry/convex_hull_trick.html
https://oi-wiki.org/dp/opt/slope/

See also: dp.go (convexHullTrick, convexHullTrick2) for DP-oriented variants.

Compared with Li Chao Tree (also in this package):
  - Li Chao Tree: arbitrary insert order, arbitrary query order, O(log n) each
  - CHT (this file): lines must be added in monotone slope order
      * chtDeque:  non-decreasing slopes + non-decreasing x queries → O(1) amortized
      * chtStack:  non-increasing slopes + arbitrary x queries → O(log n) per query

Common DP pattern:
  dp[i] = max over j<i of { a[j]*b[i] + c[j] }
  where a[j] is the "slope" and c[j] is the "intercept" of line j.

Related problems:
https://codeforces.com/problemset/problem/311/E  2500
https://codeforces.com/problemset/problem/319/C  2400
LC1235 https://leetcode.cn/problems/maximum-profit-in-job-scheduling/
https://atcoder.jp/contests/dp/tasks/dp_z
*/

type chtLine struct{ k, b int }

func (l chtLine) eval(x int) int { return l.k*x + l.b }

// chtMulCmp compares a*b with c*d using exact 128-bit signed arithmetic
// to avoid 64-bit integer overflow when slopes and intercepts are large.
func chtMulCmp(a, b, c, d int) int {
	mulSigned := func(x, y int64) (int64, uint64) {
		neg := (x < 0) != (y < 0)
		var ux, uy uint64
		if x < 0 {
			ux = uint64(-x)
		} else {
			ux = uint64(x)
		}
		if y < 0 {
			uy = uint64(-y)
		} else {
			uy = uint64(y)
		}
		uhi, ulo := bits.Mul64(ux, uy)
		if neg {
			ulo = ^ulo + 1
			uhi = ^uhi
			if ulo == 0 {
				uhi++
			}
		}
		return int64(uhi), ulo
	}
	hi1, lo1 := mulSigned(int64(a), int64(b))
	hi2, lo2 := mulSigned(int64(c), int64(d))
	if hi1 != hi2 {
		if hi1 < hi2 {
			return -1
		}
		return 1
	}
	if lo1 < lo2 {
		return -1
	}
	if lo1 > lo2 {
		return 1
	}
	return 0
}

// chtBadInc returns true if line b is never the maximum among a, b, c
// when slopes are strictly increasing (a.k < b.k < c.k).
// Equivalent to intersect(a, c) <= intersect(a, b).
func chtBadInc(a, b, c chtLine) bool {
	return chtMulCmp(a.b-c.b, b.k-a.k, a.b-b.b, c.k-a.k) <= 0
}

// chtBadDec returns true if line b is never the maximum among a, b, c
// when slopes are strictly decreasing (a.k > b.k > c.k).
// Equivalent to intersect(a, c) >= intersect(a, b).
func chtBadDec(a, b, c chtLine) bool {
	return chtMulCmp(c.b-a.b, a.k-b.k, b.b-a.b, a.k-c.k) >= 0
}

// chtDeque supports:
//
//	AddLine:  lines added in NON-DECREASING k order
//	QueryMax: x must be NON-DECREASING across queries
//
// Both operations are O(1) amortized, even when interleaved.
type chtDeque struct {
	lines []chtLine
}

func (c *chtDeque) AddLine(k, b int) {
	if n := len(c.lines); n > 0 && c.lines[n-1].k == k {
		if b <= c.lines[n-1].b {
			return
		}
		c.lines = c.lines[:n-1]
	}
	l := chtLine{k, b}
	for len(c.lines) >= 2 && chtBadInc(c.lines[len(c.lines)-2], c.lines[len(c.lines)-1], l) {
		c.lines = c.lines[:len(c.lines)-1]
	}
	c.lines = append(c.lines, l)
}

// QueryMax returns max y at x. x must be >= all previous query x values.
func (c *chtDeque) QueryMax(x int) int {
	if len(c.lines) == 0 {
		panic("chtDeque: no lines added")
	}
	for len(c.lines) >= 2 && c.lines[0].eval(x) <= c.lines[1].eval(x) {
		c.lines = c.lines[1:]
	}
	return c.lines[0].eval(x)
}

// chtStack supports:
//
//	AddLine:  lines added in NON-INCREASING k order
//	QueryMax: arbitrary x (binary search)
type chtStack struct {
	lines []chtLine
}

func (c *chtStack) AddLine(k, b int) {
	if n := len(c.lines); n > 0 && c.lines[n-1].k == k {
		if b <= c.lines[n-1].b {
			return
		}
		c.lines = c.lines[:n-1]
	}
	l := chtLine{k, b}
	for len(c.lines) >= 2 && chtBadDec(c.lines[len(c.lines)-2], c.lines[len(c.lines)-1], l) {
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
	// binary search: find first line where lines[i] >= lines[i+1] at x
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

// chtMinDeque supports minimum queries (e.g. AtCoder dp_z Frog 3, CF 319C):
//
//	AddLine:  lines added in NON-INCREASING k order
//	QueryMin: x must be NON-DECREASING across queries
//
// Internally negates (k, b) so non-increasing k becomes non-decreasing -k in chtDeque.
type chtMinDeque struct {
	inner chtDeque
}

func (c *chtMinDeque) AddLine(k, b int) { c.inner.AddLine(-k, -b) }
func (c *chtMinDeque) QueryMin(x int) int {
	return -c.inner.QueryMax(x)
}

// chtMinStack supports minimum queries with arbitrary query x:
//
//	AddLine:  lines added in NON-DECREASING k order
//	QueryMin: arbitrary x (O(log n) binary search)
type chtMinStack struct {
	inner chtStack
}

func (c *chtMinStack) AddLine(k, b int) { c.inner.AddLine(-k, -b) }
func (c *chtMinStack) QueryMin(x int) int {
	return -c.inner.QueryMax(x)
}
