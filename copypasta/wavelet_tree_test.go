package copypasta

import (
	"math"
	"math/rand"
	"sort"
	"testing"
)

func TestWaveletTreeKth(t *testing.T) {
	a := []int{3, 1, 4, 1, 5, 9, 2, 6}
	orig := append([]int(nil), a...)
	wt := newWaveletTree(a)
	// Ensure input slice was not mutated.
	for i := range a {
		if a[i] != orig[i] {
			t.Fatalf("input mutated at %d", i)
		}
	}

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

func TestWaveletTreeCountLEAndLargeValues(t *testing.T) {
	a := []int{3, 1, 4, 1, 5, 9, 2, 6}
	wt := newWaveletTree(a)

	if got := wt.CountLE(0, 7, 4); got != 5 {
		t.Fatalf("CountLE(0,7,4): want 5, got %d", got)
	}
	if got := wt.CountLE(0, 3, 1); got != 2 {
		t.Fatalf("CountLE(0,3,1): want 2, got %d", got)
	}
	if got := wt.CountLE(0, 7, 9); got != 8 {
		t.Fatalf("CountLE(0,7,9): want 8, got %d", got)
	}
	if got := wt.CountLE(0, 7, 0); got != 0 {
		t.Fatalf("CountLE(0,7,0): want 0, got %d", got)
	}
	if got := wt.CountLE(0, 7, math.MaxInt); got != 8 {
		t.Fatalf("CountLE(MaxInt): want 8, got %d", got)
	}

	// Large & negative coordinates [-1e15, 1e15] that would blow up an uncompressed wavelet tree
	rng := rand.New(rand.NewSource(88))
	bigA := make([]int, 80)
	for i := range bigA {
		bigA[i] = (rng.Intn(21) - 10) * 100_000_000_000_000
	}
	bigWT := newWaveletTree(bigA)
	for q := 0; q < 150; q++ {
		l, r := rng.Intn(len(bigA)), rng.Intn(len(bigA))
		if l > r {
			l, r = r, l
		}
		sub := append([]int(nil), bigA[l:r+1]...)
		sort.Ints(sub)
		k := 1 + rng.Intn(len(sub))
		if got := bigWT.Kth(l, r, k); got != sub[k-1] {
			t.Fatalf("bigWT Kth(%d,%d,%d): got %d, want %d", l, r, k, got, sub[k-1])
		}

		loV := (rng.Intn(25) - 12) * 100_000_000_000_000
		hiV := (rng.Intn(25) - 12) * 100_000_000_000_000
		if loV > hiV {
			loV, hiV = hiV, loV
		}
		wantRange := 0
		for _, v := range sub {
			if loV <= v && v <= hiV {
				wantRange++
			}
		}
		if got := bigWT.CountRange(l, r, loV, hiV); got != wantRange {
			t.Fatalf("bigWT CountRange(%d,%d,%d,%d): got %d, want %d", l, r, loV, hiV, got, wantRange)
		}
	}
}
