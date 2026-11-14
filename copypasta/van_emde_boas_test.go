package copypasta

import (
	"math/rand"
	"sort"
	"testing"
)

func TestVEBBasic(t *testing.T) {
	v := newVEB(16)

	if !v.IsEmpty() {
		t.Fatal("new vEB must be empty")
	}

	for _, x := range []int{3, 7, 1, 15, 0, 9} {
		v.Insert(x)
	}

	for _, x := range []int{3, 7, 1, 15, 0, 9} {
		if !v.Member(x) {
			t.Fatalf("Member(%d) should be true", x)
		}
	}
	for _, x := range []int{2, 4, 5, 6, 8, 10, 11, 12, 13, 14} {
		if v.Member(x) {
			t.Fatalf("Member(%d) should be false", x)
		}
	}

	if v.Min() != 0 {
		t.Fatalf("Min: want 0, got %d", v.Min())
	}
	if v.Max() != 15 {
		t.Fatalf("Max: want 15, got %d", v.Max())
	}
}

func TestVEBSuccessorPredecessor(t *testing.T) {
	v := newVEB(32)
	elems := []int{2, 5, 11, 17, 23, 29}
	for _, x := range elems {
		v.Insert(x)
	}

	// Successor
	cases := [][2]int{{1, 2}, {2, 5}, {5, 11}, {11, 17}, {17, 23}, {23, 29}, {29, vebNil}, {30, vebNil}}
	for _, c := range cases {
		if got := v.Successor(c[0]); got != c[1] {
			t.Errorf("Successor(%d): want %d, got %d", c[0], c[1], got)
		}
	}

	// Predecessor
	cases = [][2]int{{30, 29}, {29, 23}, {23, 17}, {17, 11}, {11, 5}, {5, 2}, {2, vebNil}, {1, vebNil}}
	for _, c := range cases {
		if got := v.Predecessor(c[0]); got != c[1] {
			t.Errorf("Predecessor(%d): want %d, got %d", c[0], c[1], got)
		}
	}
}

func TestVEBDelete(t *testing.T) {
	v := newVEB(16)
	for _, x := range []int{1, 3, 5, 7, 9} {
		v.Insert(x)
	}

	v.Delete(5)
	if v.Member(5) {
		t.Fatal("5 should be deleted")
	}
	if v.Successor(3) != 7 {
		t.Fatalf("Successor(3) after deleting 5: want 7, got %d", v.Successor(3))
	}

	// delete min
	v.Delete(1)
	if v.Min() != 3 {
		t.Fatalf("Min after deleting 1: want 3, got %d", v.Min())
	}

	// delete max
	v.Delete(9)
	if v.Max() != 7 {
		t.Fatalf("Max after deleting 9: want 7, got %d", v.Max())
	}
}

func TestVEBRandomized(t *testing.T) {
	const U = 1 << 10
	v := newVEB(U)
	present := make(map[int]bool)

	ops := rand.Perm(U)[:200]
	for _, x := range ops {
		v.Insert(x)
		present[x] = true
	}

	sorted := make([]int, 0, len(present))
	for x := range present {
		sorted = append(sorted, x)
	}
	sort.Ints(sorted)

	// verify successor chain
	cur := v.Min()
	for _, want := range sorted {
		if cur != want {
			t.Fatalf("successor chain mismatch: got %d want %d", cur, want)
		}
		cur = v.Successor(cur)
	}
	if cur != vebNil {
		t.Fatalf("chain should end with vebNil, got %d", cur)
	}
}

func BenchmarkVEBInsert(b *testing.B) {
	v := newVEB(1 << 20)
	b.ResetTimer()
	for i := range b.N {
		v.Insert(rand.Intn(1 << 20))
		_ = i
	}
}

func BenchmarkVEBSuccessor(b *testing.B) {
	v := newVEB(1 << 20)
	for i := range 10000 {
		v.Insert(rand.Intn(1 << 20))
		_ = i
	}
	b.ResetTimer()
	for range b.N {
		v.Successor(rand.Intn(1 << 20))
	}
}
