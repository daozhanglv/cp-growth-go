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

func TestVEBDuplicateInsertAndMissingDelete(t *testing.T) {
	v := newVEB(16)

	// Delete on empty tree should be a no-op.
	v.Delete(0)
	v.Delete(5)
	if !v.IsEmpty() {
		t.Fatal("tree should still be empty")
	}

	// Single element: deleting a different value must not clear the set.
	v.Insert(5)
	v.Delete(3)
	if !v.Member(5) || v.Min() != 5 || v.Max() != 5 {
		t.Fatalf("deleting missing element corrupted single-element tree: min=%d max=%d", v.Min(), v.Max())
	}

	// Duplicate insert of min and non-min must not corrupt invariants.
	v.Insert(5)
	v.Insert(1)
	v.Insert(1) // duplicate min
	v.Insert(9)
	v.Insert(9) // duplicate max
	v.Delete(1)
	if v.Member(1) || v.Min() != 5 {
		t.Fatalf("after deleting duplicate-inserted min 1: member=%v min=%d", v.Member(1), v.Min())
	}
	v.Delete(7) // missing element in multi-element tree
	if !v.Member(5) || !v.Member(9) || v.Min() != 5 || v.Max() != 9 {
		t.Fatalf("deleting missing element corrupted tree: min=%d max=%d", v.Min(), v.Max())
	}

	// Out-of-range operations should be safe no-ops and return valid bounds.
	v.Insert(-1)
	v.Insert(16)
	v.Delete(-1)
	v.Delete(16)
	if v.Member(-1) || v.Member(16) {
		t.Fatal("out-of-range elements should never be members")
	}
	// Remaining set is {5, 9}
	if got := v.Successor(-10); got != 5 {
		t.Fatalf("Successor(-10): want 5, got %d", got)
	}
	if got := v.Successor(16); got != vebNil {
		t.Fatalf("Successor(16): want vebNil, got %d", got)
	}
	if got := v.Successor(100); got != vebNil {
		t.Fatalf("Successor(100): want vebNil, got %d", got)
	}
	if got := v.Predecessor(-10); got != vebNil {
		t.Fatalf("Predecessor(-10): want vebNil, got %d", got)
	}
	if got := v.Predecessor(16); got != 9 {
		t.Fatalf("Predecessor(16): want 9, got %d", got)
	}
	if got := v.Predecessor(100); got != 9 {
		t.Fatalf("Predecessor(100): want 9, got %d", got)
	}

	// Small and degenerate universes u = 0, 1, 2, 3
	emptyU := newVEB(0)
	emptyU.Insert(0)
	if !emptyU.IsEmpty() || emptyU.Successor(-1) != vebNil || emptyU.Predecessor(1) != vebNil {
		t.Fatal("u=0 vEB must stay empty")
	}
	for _, u := range []int{1, 2, 3, 4} {
		small := newVEB(u)
		small.Insert(0)
		if u > 1 {
			small.Insert(u - 1)
		}
		if got := small.Successor(-1); got != 0 {
			t.Fatalf("u=%d Successor(-1): want 0, got %d", u, got)
		}
		if got := small.Successor(u + 5); got != vebNil {
			t.Fatalf("u=%d Successor(%d): want vebNil, got %d", u, u+5, got)
		}
		if got := small.Predecessor(-1); got != vebNil {
			t.Fatalf("u=%d Predecessor(-1): want vebNil, got %d", u, got)
		}
		if got := small.Predecessor(u + 5); got != u-1 {
			t.Fatalf("u=%d Predecessor(%d): want %d, got %d", u, u+5, u-1, got)
		}
	}
}

func TestVEBRandomizedBruteForce(t *testing.T) {
	const U = 1 << 8
	rng := rand.New(rand.NewSource(12345))
	v := newVEB(U)
	set := make(map[int]bool)

	bruteMin := func() int {
		res := vebNil
		for x := range set {
			if res == vebNil || x < res {
				res = x
			}
		}
		return res
	}
	bruteMax := func() int {
		res := vebNil
		for x := range set {
			if res == vebNil || x > res {
				res = x
			}
		}
		return res
	}
	bruteSucc := func(q int) int {
		res := vebNil
		for x := range set {
			if x > q && (res == vebNil || x < res) {
				res = x
			}
		}
		return res
	}
	brutePred := func(q int) int {
		res := vebNil
		for x := range set {
			if x < q && (res == vebNil || x > res) {
				res = x
			}
		}
		return res
	}

	for step := 0; step < 3000; step++ {
		x := rng.Intn(U)
		switch rng.Intn(4) {
		case 0, 1: // insert (including duplicates)
			v.Insert(x)
			set[x] = true
		case 2: // delete (including missing elements)
			v.Delete(x)
			delete(set, x)
		case 3: // queries
			if got, want := v.Member(x), set[x]; got != want {
				t.Fatalf("step %d Member(%d): got %v want %v", step, x, got, want)
			}
			if got, want := v.Min(), bruteMin(); got != want {
				t.Fatalf("step %d Min(): got %d want %d", step, got, want)
			}
			if got, want := v.Max(), bruteMax(); got != want {
				t.Fatalf("step %d Max(): got %d want %d", step, got, want)
			}
			if got, want := v.Successor(x), bruteSucc(x); got != want {
				t.Fatalf("step %d Successor(%d): got %d want %d", step, x, got, want)
			}
			if got, want := v.Predecessor(x), brutePred(x); got != want {
				t.Fatalf("step %d Predecessor(%d): got %d want %d", step, x, got, want)
			}
		}
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
