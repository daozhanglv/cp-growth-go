package copypasta

import (
	"sort"
	"testing"
)

func TestWaveletTreeKth(t *testing.T) {
	a := []int{3, 1, 4, 1, 5, 9, 2, 6}
	wt := newWaveletTree(append([]int(nil), a...), 1, 9)

	// full array sorted: [1,1,2,3,4,5,6,9]
	sorted := append([]int(nil), a...)
	sort.Ints(sorted)
	for k := 1; k <= len(a); k++ {
		if got := wt.Kth(0, len(a)-1, k); got != sorted[k-1] {
			t.Fatalf("Kth(0,%d,%d): want %d, got %d", len(a)-1, k, sorted[k-1], got)
		}
	}

	// subarray [1..4] = [1,4,1,5] → sorted [1,1,4,5]
	sub := []int{1, 4, 1, 5}
	sort.Ints(sub)
	for k := 1; k <= 4; k++ {
		if got := wt.Kth(1, 4, k); got != sub[k-1] {
			t.Fatalf("Kth(1,4,%d): want %d, got %d", k, sub[k-1], got)
		}
	}
}

func TestWaveletTreeCountLE(t *testing.T) {
	a := []int{3, 1, 4, 1, 5, 9, 2, 6}
	wt := newWaveletTree(append([]int(nil), a...), 1, 9)

	// count elements ≤ 4 in full array: [3,1,4,1,2] = 5
	if got := wt.CountLE(0, 7, 4); got != 5 {
		t.Fatalf("CountLE(0,7,4): want 5, got %d", got)
	}
	// count elements ≤ 1 in [0..3]: [1,1] = 2
	if got := wt.CountLE(0, 3, 1); got != 2 {
		t.Fatalf("CountLE(0,3,1): want 2, got %d", got)
	}
	// count elements ≤ 9 in full array: all 8
	if got := wt.CountLE(0, 7, 9); got != 8 {
		t.Fatalf("CountLE(0,7,9): want 8, got %d", got)
	}
	// count elements ≤ 0: 0
	if got := wt.CountLE(0, 7, 0); got != 0 {
		t.Fatalf("CountLE(0,7,0): want 0, got %d", got)
	}
}
