package copypasta

import (
	"math/rand"
	"testing"
)

func TestBerlekampMasseyFibonacci(t *testing.T) {
	mod := int64(1e9 + 7)
	// Fibonacci: 0,1,1,2,3,5,8,13,21,34
	s := []int64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	rec := berlekampMasseyMod(s, mod)
	// Expected recurrence: s[i] = 1*s[i-1] + 1*s[i-2]
	if len(rec) != 2 || rec[0] != 1 || rec[1] != 1 {
		t.Fatalf("Fibonacci recurrence: want [1 1], got %v", rec)
	}
}

func TestBerlekampMasseyGeometric(t *testing.T) {
	mod := int64(1e9 + 7)
	// Geometric: 1,2,4,8,16,32 (ratio 2)
	s := []int64{1, 2, 4, 8, 16, 32}
	rec := berlekampMasseyMod(s, mod)
	// s[i] = 2*s[i-1]
	if len(rec) != 1 || rec[0] != 2 {
		t.Fatalf("geometric recurrence: want [2], got %v", rec)
	}
}

func TestBerlekampMasseyConstant(t *testing.T) {
	mod := int64(1e9 + 7)
	// Constant sequence 5,5,5,5,5
	s := []int64{5, 5, 5, 5, 5}
	rec := berlekampMasseyMod(s, mod)
	if len(rec) != 1 || rec[0] != 1 {
		t.Fatalf("constant recurrence: want [1], got %v", rec)
	}
}

func TestBerlekampMasseyUnreducedAndNegativeInputs(t *testing.T) {
	const p int64 = 1_000_000_007
	// Fibonacci shifted by large multiples of p (around +/- 5e18) that would overflow C[j]*s[i-j] without normalization.
	fib := []int64{0, 1, 1, 2, 3, 5, 8, 13, 21, 34}
	noisy := make([]int64, len(fib))
	for i, v := range fib {
		if i%2 == 0 {
			noisy[i] = v + 5_000_000_000*p
		} else {
			noisy[i] = v - 5_000_000_000*p
		}
	}
	rec := berlekampMasseyMod(noisy, p)
	if len(rec) != 2 || rec[0] != 1 || rec[1] != 1 {
		t.Fatalf("noisy Fibonacci recurrence: want [1 1], got %v", rec)
	}
}

func TestKitamasaModAndGuessNth(t *testing.T) {
	const p int64 = 998244353
	rng := rand.New(rand.NewSource(77))
	for tc := 0; tc < 30; tc++ {
		L := 1 + rng.Intn(6)
		trueRec := make([]int64, L)
		for i := range trueRec {
			trueRec[i] = int64(rng.Intn(int(p)))
		}
		terms := make([]int64, 2*L+50)
		for i := 0; i < L; i++ {
			terms[i] = int64(rng.Intn(int(p)))
		}
		for i := L; i < len(terms); i++ {
			var sum int64
			for j := 0; j < L; j++ {
				sum = (sum + trueRec[j]*terms[i-1-j]) % p
			}
			terms[i] = sum
		}

		prefix := terms[:2*L+2]
		for k := 0; k < len(terms); k++ {
			if got := guessNthMod(prefix, int64(k), p); got != terms[k] {
				t.Fatalf("tc %d k=%d: guessNthMod got %d, want %d", tc, k, got, terms[k])
			}
		}
	}

	// Also verify math_matrix.go's int-based guessNth (modulo global mod = 1e9+7).
	fibInt := []int{0, 1, 1, 2, 3, 5, 8, 13}
	a, b := 0, 1
	for n := 0; n <= 60; n++ {
		if got := guessNth(fibInt, n); got != a {
			t.Fatalf("guessNth(fib, %d): got %d, want %d", n, got, a)
		}
		a, b = b, (a+b)%mod
	}
}
