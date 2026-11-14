package copypasta

import (
	"sort"
	"testing"
)

func TestSkewHeapBasic(t *testing.T) {
	var h *skewNode
	vals := []int{5, 3, 7, 1, 9, 2, 8, 4, 6}
	for _, v := range vals {
		h = skewPush(h, v)
	}

	sorted := append([]int(nil), vals...)
	sort.Ints(sorted)

	for _, want := range sorted {
		if h == nil {
			t.Fatal("heap is empty too soon")
		}
		if h.val != want {
			t.Fatalf("min: want %d, got %d", want, h.val)
		}
		h = skewPop(h)
	}
	if h != nil {
		t.Fatal("heap should be empty")
	}
}

func TestSkewHeapMerge(t *testing.T) {
	var h1, h2 *skewNode
	for _, v := range []int{2, 4, 6} {
		h1 = skewPush(h1, v)
	}
	for _, v := range []int{1, 3, 5} {
		h2 = skewPush(h2, v)
	}
	h := skewMerge(h1, h2)

	for _, want := range []int{1, 2, 3, 4, 5, 6} {
		if h == nil || h.val != want {
			t.Fatalf("merge order: want %d, got %v", want, h)
		}
		h = skewPop(h)
	}
}
