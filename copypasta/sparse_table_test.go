package copypasta

import (
	"math/rand"
	"testing"
)

func gcdInt(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func TestSparseTableMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(801))
	for n := 1; n <= 45; n++ {
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(100) + 1
		}
		stMax := newSparseTable(a, func(x, y int) int { return max(x, y) })
		stMin := newSparseTable(a, func(x, y int) int { return min(x, y) })
		stGCD := newSparseTable(a, gcdInt)
		stIdx := newSparseTableWithIndex(a)

		for l := 0; l < n; l++ {
			curMax, curMin, curGCD := a[l], a[l], a[l]
			minIdx := l
			for r := l + 1; r <= n; r++ {
				v := a[r-1]
				if v > curMax {
					curMax = v
				}
				if v < curMin {
					curMin = v
					minIdx = r - 1
				}
				curGCD = gcdInt(curGCD, v)

				if got := stMax.query(l, r); got != curMax {
					t.Fatalf("n=%d max[%d,%d): got %d, want %d", n, l, r, got, curMax)
				}
				if got := stMin.query(l, r); got != curMin {
					t.Fatalf("n=%d min[%d,%d): got %d, want %d", n, l, r, got, curMin)
				}
				if got := stGCD.query(l, r); got != curGCD {
					t.Fatalf("n=%d gcd[%d,%d): got %d, want %d", n, l, r, got, curGCD)
				}
				if got := stIdx.query(l, r); got != minIdx {
					t.Fatalf("n=%d minIdx[%d,%d): got %d, want %d", n, l, r, got, minIdx)
				}
			}
		}
	}
}
