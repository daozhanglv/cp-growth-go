package copypasta

/*
Van Emde Boas Tree (vEB Tree)
https://en.wikipedia.org/wiki/Van_Emde_Boas_tree
https://www.luogu.com.cn/blog/RuntimeErrror/ni-suo-bu-zhi-dao-de-shu-ju-jie-gou-van-emde-boas-shu

Maintains a dynamic set S ⊆ [0, U-1] and supports all operations in O(log log U).
This beats O(log n) of BSTs and O(log U) of segment trees for predecessor/successor.

Time complexity:  O(log log U) per operation (insert, delete, member, succ, pred, min, max)
Space complexity: O(U) — pre-allocates the full universe; not suitable for large U

Suitable for competitive programming when:
  - U ≤ 2^20 (≈ 1 million) — memory is reasonable
  - Many predecessor/successor queries on a sparse set
  - Faster than segment tree for this specific pattern

The tree recursively splits the universe into √U clusters.
A key x is stored in cluster[high(x)] at position low(x).
The min of each node is NOT stored recursively (lazy insertion) to achieve O(log log U).

Related problems:
https://codeforces.com/blog/entry/101663  (E — vEB as a building block)
https://leetcode.cn/problems/online-majority-element-in-subarray/  (alternative use)
https://codeforces.com/problemset/problem/786/C  2400  (offline intervals, successor query)
*/

const vebNil = -1

// veb is a Van Emde Boas tree over the universe [0, u-1].
// u should be a power of 2; call newVEB(nextPow2(n)) for universe [0, n-1].
type veb struct {
	min, max int  // vebNil if empty
	u        int  // universe size
	sqrtU    int  // lower-sqrt of u = 2^(floor(log2(u)/2)); cluster size
	summary  *veb // vEB over cluster indices
	cluster  []*veb
}

// newVEB creates a vEB tree for the universe [0, u-1].
// u is rounded up to the next power of 2 automatically.
func newVEB(u int) *veb {
	// round up to power of 2
	pow := 1
	for pow < u {
		pow <<= 1
	}
	return buildVEB(pow)
}

func buildVEB(u int) *veb {
	v := &veb{min: vebNil, max: vebNil, u: u}
	if u <= 2 {
		v.sqrtU = 1
		return v
	}
	// bits = log2(u); sqrtU = 2^(bits/2)
	bits := 0
	for (1 << bits) < u {
		bits++
	}
	v.sqrtU = 1 << (bits / 2)
	numClusters := u / v.sqrtU // == 2^(bits - bits/2) = upper-sqrt

	v.summary = buildVEB(numClusters)
	v.cluster = make([]*veb, numClusters)
	for i := range v.cluster {
		v.cluster[i] = buildVEB(v.sqrtU)
	}
	return v
}

func (v *veb) high(x int) int  { return x / v.sqrtU }
func (v *veb) low(x int) int   { return x % v.sqrtU }
func (v *veb) idx(h, l int) int { return h*v.sqrtU + l }

// IsEmpty reports whether the set is empty.
func (v *veb) IsEmpty() bool { return v.min == vebNil }

// Min returns the minimum element, or vebNil if empty.
func (v *veb) Min() int { return v.min }

// Max returns the maximum element, or vebNil if empty.
func (v *veb) Max() int { return v.max }

// Member reports whether x is in the set.
func (v *veb) Member(x int) bool {
	switch x {
	case v.min, v.max:
		return true
	}
	if v.u <= 2 {
		return false
	}
	return v.cluster[v.high(x)].Member(v.low(x))
}

// Insert adds x to the set. No-op if already present.
func (v *veb) Insert(x int) {
	if v.min == vebNil {
		v.min, v.max = x, x
		return
	}
	if x < v.min {
		x, v.min = v.min, x // push old min down
	}
	if v.u > 2 {
		h, l := v.high(x), v.low(x)
		if v.cluster[h].IsEmpty() {
			v.summary.Insert(h)
		}
		v.cluster[h].Insert(l)
	}
	if x > v.max {
		v.max = x
	}
}

// Delete removes x from the set. No-op if not present.
func (v *veb) Delete(x int) {
	if v.min == v.max {
		// only one element
		v.min, v.max = vebNil, vebNil
		return
	}
	if v.u <= 2 {
		// universe {0,1} with two elements: remove x, keep the other
		if x == 0 {
			v.min = 1
		} else {
			v.min = 0
		}
		v.max = v.min
		return
	}
	if x == v.min {
		// promote the smallest element from clusters to become the new min
		firstCluster := v.summary.min
		x = v.idx(firstCluster, v.cluster[firstCluster].min)
		v.min = x // update min without recursing further
	}
	h, l := v.high(x), v.low(x)
	v.cluster[h].Delete(l)
	if v.cluster[h].IsEmpty() {
		v.summary.Delete(h)
		if x == v.max {
			summaryMax := v.summary.max
			if summaryMax == vebNil {
				v.max = v.min
			} else {
				v.max = v.idx(summaryMax, v.cluster[summaryMax].max)
			}
		}
	} else if x == v.max {
		v.max = v.idx(h, v.cluster[h].max)
	}
}

// Successor returns the smallest element > x, or vebNil if none.
func (v *veb) Successor(x int) int {
	if v.u <= 2 {
		if x == 0 && v.max == 1 {
			return 1
		}
		return vebNil
	}
	if !v.IsEmpty() && x < v.min {
		return v.min
	}
	h, l := v.high(x), v.low(x)
	if !v.cluster[h].IsEmpty() && l < v.cluster[h].max {
		return v.idx(h, v.cluster[h].Successor(l))
	}
	succCluster := v.summary.Successor(h)
	if succCluster == vebNil {
		return vebNil
	}
	return v.idx(succCluster, v.cluster[succCluster].min)
}

// Predecessor returns the largest element < x, or vebNil if none.
func (v *veb) Predecessor(x int) int {
	if v.u <= 2 {
		if x == 1 && v.min == 0 {
			return 0
		}
		return vebNil
	}
	if !v.IsEmpty() && x > v.max {
		return v.max
	}
	h, l := v.high(x), v.low(x)
	if !v.cluster[h].IsEmpty() && l > v.cluster[h].min {
		return v.idx(h, v.cluster[h].Predecessor(l))
	}
	predCluster := v.summary.Predecessor(h)
	if predCluster == vebNil {
		// x > v.min but no predecessor in clusters: v.min is the predecessor
		if v.min != vebNil && x > v.min {
			return v.min
		}
		return vebNil
	}
	return v.idx(predCluster, v.cluster[predCluster].max)
}
