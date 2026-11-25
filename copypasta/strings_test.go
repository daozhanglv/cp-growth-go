package copypasta

import (
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestKMPZAndManacherMatchBruteForce(t *testing.T) {
	if got := kmpCalcPi(""); len(got) != 0 {
		t.Fatalf("kmpCalcPi(\"\"): got %v", got)
	}
	if got := calcZ(""); len(got) != 0 {
		t.Fatalf("calcZ(\"\"): got %v", got)
	}
	if p, c := calcMinPeriod(""); p != "" || c != 0 {
		t.Fatalf("calcMinPeriod(\"\"): got (%q, %d)", p, c)
	}
	if got := countPalindromesManacher(""); got != 0 {
		t.Fatalf("countPalindromesManacher(\"\"): got %d", got)
	}

	// Regression cases for zCompare: suffix longer than pattern & suffix proper prefix of pattern.
	if got := zCompare("abc", "ab"); !slices.Equal(got, []int{1, 1, 1}) {
		t.Fatalf("zCompare(\"abc\", \"ab\"): got %v, want [1 1 1]", got)
	}
	if got := zCompare("a", "ab"); !slices.Equal(got, []int{-1}) {
		t.Fatalf("zCompare(\"a\", \"ab\"): got %v, want [-1]", got)
	}

	rng := rand.New(rand.NewSource(909))
	for tc := 0; tc < 60; tc++ {
		tb := make([]byte, 1+rng.Intn(45))
		for i := range tb {
			tb[i] = byte('a' + rng.Intn(3))
		}
		text := string(tb)

		pb := make([]byte, 1+rng.Intn(8))
		for i := range pb {
			pb[i] = byte('a' + rng.Intn(3))
		}
		pattern := string(pb)

		// 1. KMP prefix function vs brute-force border
		pi := kmpCalcPi(text)
		for i := range text {
			sub := text[:i+1]
			want := 0
			for l := 1; l < len(sub); l++ {
				if sub[:l] == sub[len(sub)-l:] {
					want = l
				}
			}
			if pi[i] != want {
				t.Fatalf("tc %d pi[%d] on %q: got %d, want %d", tc, i, text, pi[i], want)
			}
		}

		// 2. kmpSearch and zSearch vs brute force
		var wantPos []int
		for i := 0; i+len(pattern) <= len(text); i++ {
			if text[i:i+len(pattern)] == pattern {
				wantPos = append(wantPos, i)
			}
		}
		if got := kmpSearch(text, pattern); !slices.Equal(got, wantPos) {
			t.Fatalf("tc %d kmpSearch(%q, %q): got %v, want %v", tc, text, pattern, got, wantPos)
		}
		if got := zSearch(text, pattern); !slices.Equal(got, wantPos) {
			t.Fatalf("tc %d zSearch(%q, %q): got %v, want %v", tc, text, pattern, got, wantPos)
		}

		// 3. calcMinPeriod vs brute force
		unit, cnt := calcMinPeriod(text)
		if strings.Repeat(unit, cnt) != text {
			t.Fatalf("tc %d calcMinPeriod(%q) = (%q, %d) does not reconstruct", tc, text, unit, cnt)
		}
		for l := 1; l < len(unit); l++ {
			if len(text)%l == 0 && strings.Repeat(text[:l], len(text)/l) == text {
				t.Fatalf("tc %d calcMinPeriod(%q) missed shorter period %q", tc, text, text[:l])
			}
		}

		// 4. calcZ & zCompare vs brute force
		z := calcZ(text)
		for i := range text {
			h := 0
			for i+h < len(text) && text[h] == text[i+h] {
				h++
			}
			if z[i] != h {
				t.Fatalf("tc %d z[%d] on %q: got %d, want %d", tc, i, text, z[i], h)
			}
		}
		zc := zCompare(text, pattern)
		for i := range text {
			if want := strings.Compare(text[i:], pattern); zc[i] != want {
				t.Fatalf("tc %d zCompare[%d] (%q vs %q): got %d, want %d", tc, i, text[i:], pattern, zc[i], want)
			}
		}

		// 5. Manacher palindrome range checks & count vs brute force
		hl := manacherHalfLen(text)
		wantPalCnt := 0
		for l := 0; l < len(text); l++ {
			for r := l; r < len(text); r++ {
				isPal := true
				for p, q := l, r; p < q; p, q = p+1, q-1 {
					if text[p] != text[q] {
						isPal = false
						break
					}
				}
				if isPal {
					wantPalCnt++
				}
				if got := hl[l+r+2] > r-l+1; got != isPal {
					t.Fatalf("tc %d manacher[%d..%d] on %q: got %v, want %v", tc, l, r, text, got, isPal)
				}
			}
		}
		if got := countPalindromesManacher(text); got != wantPalCnt {
			t.Fatalf("tc %d countPalindromesManacher(%q): got %d, want %d", tc, text, got, wantPalCnt)
		}
	}
}
