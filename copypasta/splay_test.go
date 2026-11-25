package copypasta

import (
	"math/rand"
	"slices"
	"sort"
	"testing"
)

func TestSplayMultisetMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(601))
	sp := newSplay()
	var arr []int // sorted multiset

	// Note: newSplay() inserts two sentinels (-2e9 and +2e9), so rank(k) is 1 + count(<= k)
	// and kth user element (1-indexed) is at splay index k+1.
	for step := 0; step < 400; step++ {
		key := rng.Intn(40) - 20
		switch rng.Intn(3) {
		case 0, 1: // insert 1 occurrence
			sp.put(spKeyType(key), 1)
			idx := sort.SearchInts(arr, key)
			arr = slices.Insert(arr, idx, key)
		case 2: // delete 1 occurrence if present
			sp.delete(spKeyType(key))
			if idx := sort.SearchInts(arr, key); idx < len(arr) && arr[idx] == key {
				arr = slices.Delete(arr, idx, idx+1)
			}
		}

		// Verify total size (2 sentinels + unique keys? Wait: spNode stores multiplicity in val, sz counts distinct nodes!)
		// Let's verify in-order traversal of keys & multiplicities matches arr!
		var flat []int
		var dfs func(o *spNode)
		dfs = func(o *spNode) {
			if o == nil {
				return
			}
			dfs(o.lr[0])
			if o.key != splayMin && o.key != splayMax {
				for c := spValueType(0); c < o.val; c++ {
					flat = append(flat, int(o.key))
				}
			}
			dfs(o.lr[1])
		}
		dfs(sp.root)
		if !slices.Equal(flat, arr) {
			t.Fatalf("step %d multiset mismatch:\n got %v\nwant %v", step, flat, arr)
		}
	}
}

func TestSplaySplitMergeRangeFlip(t *testing.T) {
	rng := rand.New(rand.NewSource(602))
	const n = 25
	root := buildSplay(1, n)
	arr := make([]int, n)
	for i := range arr {
		arr[i] = i + 1
	}

	reverseSubrange := func(root *spNode, l, r int) *spNode {
		// 1-indexed [l, r]
		var left, mid, right *spNode
		if l > 1 {
			left, mid = root.split(l - 1)
		} else {
			mid = root
		}
		mid, right = mid.split(r - l + 1)
		mid.apply()
		if left != nil {
			root = left.merge(mid)
		} else {
			root = mid
		}
		if right != nil {
			root = root.merge(right)
		}
		return root
	}

	for step := 0; step < 150; step++ {
		l, r := 1+rng.Intn(n), 1+rng.Intn(n)
		if l > r {
			l, r = r, l
		}
		root = reverseSubrange(root, l, r)
		slices.Reverse(arr[l-1 : r])

		for k := 1; k <= n; k++ {
			root = root.splay(k)
			if int(root.key) != arr[k-1] {
				t.Fatalf("step %d kth(%d): got %d, want %d", step, k, root.key, arr[k-1])
			}
		}
	}
}
