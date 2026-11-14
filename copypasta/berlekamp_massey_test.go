package copypasta

import "testing"

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
