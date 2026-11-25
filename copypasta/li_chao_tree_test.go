package copypasta

import (
	"math"
	"math/rand"
	"testing"
)

func TestLiChaoTreeBasic(t *testing.T) {
	// lines: y=2x+1, y=-x+10, y=x+3
	lct := newLiChaoTree(0, 10)
	lct.AddLine(2, 1)
	lct.AddLine(-1, 10)
	lct.AddLine(1, 3)

	// at x=0: max(1, 10, 3) = 10
	if got := lct.QueryMax(0); got != 10 {
		t.Fatalf("QueryMax(0): want 10, got %d", got)
	}
	// at x=10: max(21, 0, 13) = 21
	if got := lct.QueryMax(10); got != 21 {
		t.Fatalf("QueryMax(10): want 21, got %d", got)
	}
	// at x=3: max(7, 7, 6) = 7
	if got := lct.QueryMax(3); got != 7 {
		t.Fatalf("QueryMax(3): want 7, got %d", got)
	}
}

func TestLiChaoTreeEmpty(t *testing.T) {
	lct := newLiChaoTree(0, 100)
	if got := lct.QueryMax(50); got != math.MinInt {
		t.Fatalf("empty tree: want MinInt, got %d", got)
	}
	minTree := newLiChaoMinTree(0, 100)
	if got := minTree.QueryMin(50); got != math.MaxInt {
		t.Fatalf("empty min tree: want MaxInt, got %d", got)
	}
}

func TestLiChaoTreeMany(t *testing.T) {
	lct := newLiChaoTree(-100, 100)
	// horizontal lines y = c for c = 0..49
	for c := 0; c < 50; c++ {
		lct.AddLine(0, c)
	}
	// max should be 49 everywhere
	for _, x := range []int{-100, -50, 0, 50, 100} {
		if got := lct.QueryMax(x); got != 49 {
			t.Fatalf("QueryMax(%d): want 49, got %d", x, got)
		}
	}
	// add steep line y = 10x + 0; at x=10 it gives 100, beating 49
	lct.AddLine(10, 0)
	if got := lct.QueryMax(10); got != 100 {
		t.Fatalf("QueryMax(10) after steep line: want 100, got %d", got)
	}
}

func TestLiChaoSegmentMergeAndMinBruteForce(t *testing.T) {
	const lo, hi = -80, 80
	rng := rand.New(rand.NewSource(2026))
	type seg struct{ l, r, k, b int }
	var segs1, segs2 []seg

	t1 := newLiChaoMinTree(lo, hi)
	t2 := newLiChaoMinTree(lo, hi)
	for i := 0; i < 60; i++ {
		l := lo + rng.Intn(hi-lo+1)
		r := lo + rng.Intn(hi-lo+1)
		if l > r {
			l, r = r, l
		}
		k := rng.Intn(41) - 20
		b := rng.Intn(401) - 200
		if i%2 == 0 {
			segs1 = append(segs1, seg{l, r, k, b})
			t1.AddSegment(l, r, k, b)
		} else {
			segs2 = append(segs2, seg{l, r, k, b})
			t2.AddSegment(l, r, k, b)
		}
	}

	// Verify t1 before merge
	for x := lo; x <= hi; x++ {
		want := math.MaxInt
		for _, s := range segs1 {
			if s.l <= x && x <= s.r {
				if v := s.k*x + s.b; v < want {
					want = v
				}
			}
		}
		if got := t1.QueryMin(x); got != want {
			t.Fatalf("before merge QueryMin(%d): got %d want %d", x, got, want)
		}
	}

	// Merge t2 into t1 and verify union of segments
	t1.Merge(t2)
	all := append(segs1, segs2...)
	for x := lo; x <= hi; x++ {
		want := math.MaxInt
		for _, s := range all {
			if s.l <= x && x <= s.r {
				if v := s.k*x + s.b; v < want {
					want = v
				}
			}
		}
		if got := t1.QueryMin(x); got != want {
			t.Fatalf("after merge QueryMin(%d): got %d want %d", x, got, want)
		}
	}
}
