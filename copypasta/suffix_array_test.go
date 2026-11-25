package copypasta

import (
	"math/rand"
	"sort"
	"testing"
)

func TestBuildSABanana(t *testing.T) {
	// "banana": suffixes in lexicographic order:
	// "a"(5), "ana"(3), "anana"(1), "banana"(0), "na"(4), "nana"(2)
	// sa = [5, 3, 1, 0, 4, 2]
	s := "banana"
	sa := buildSA(s)
	expected := []int{5, 3, 1, 0, 4, 2}
	for i, v := range expected {
		if sa[i] != v {
			t.Fatalf("sa[%d]: want %d, got %d", i, v, sa[i])
		}
	}
}

func TestBuildSASorted(t *testing.T) {
	// "abcd": each suffix starts unique, sa = [0,1,2,3]
	sa := buildSA("abcd")
	for i, v := range sa {
		if v != i {
			t.Fatalf("sa[%d]: want %d, got %d", i, i, v)
		}
	}
}

func TestBuildSAAllSame(t *testing.T) {
	// "aaaa": suffixes "aaaa","aaa","aa","a" -> sorted: "a","aa","aaa","aaaa"
	// sa = [3,2,1,0]
	sa := buildSA("aaaa")
	expected := []int{3, 2, 1, 0}
	for i, v := range expected {
		if sa[i] != v {
			t.Fatalf("sa[%d]: want %d, got %d", i, v, sa[i])
		}
	}
}

func TestBuildSAOrdering(t *testing.T) {
	// Verify SA produces lexicographically sorted suffixes
	for _, s := range []string{"mississippi", "abracadabra", "aababcabcd", "zzyzx"} {
		sa := buildSA(s)
		suffixes := make([]string, len(s))
		for i, v := range sa {
			suffixes[i] = s[v:]
		}
		if !sort.StringsAreSorted(suffixes) {
			t.Fatalf("SA not sorted for %q", s)
		}
	}
}

func TestBuildLCPBanana(t *testing.T) {
	s := "banana"
	sa := buildSA(s)
	lcp := buildLCP(s, sa)
	// sa = [5,3,1,0,4,2], lcp = [0,1,3,0,0,2]
	expected := []int{0, 1, 3, 0, 0, 2}
	for i, v := range expected {
		if lcp[i] != v {
			t.Fatalf("lcp[%d]: want %d, got %d", i, v, lcp[i])
		}
	}
}

func TestCountDistinctSubstrings(t *testing.T) {
	// "aab": distinct substrings: "a","aa","aab","ab","b" = 5
	if got := countDistinctSubstrings("aab"); got != 5 {
		t.Fatalf("countDistinctSubstrings(\"aab\"): want 5, got %d", got)
	}
	// "abc": all substrings distinct = 3+2+1 = 6
	if got := countDistinctSubstrings("abc"); got != 6 {
		t.Fatalf("countDistinctSubstrings(\"abc\"): want 6, got %d", got)
	}
}

func TestLongestCommonSubstring(t *testing.T) {
	if got := longestCommonSubstring("abcde", "cdefg"); got != 3 { // "cde"
		t.Fatalf("LCS(abcde,cdefg): want 3, got %d", got)
	}
	if got := longestCommonSubstring("abc", "xyz"); got != 0 {
		t.Fatalf("LCS(abc,xyz): want 0, got %d", got)
	}
	if got := longestCommonSubstring("abcabc", "bcabca"); got != 5 { // "bcabc" or "cabca"
		t.Fatalf("LCS(abcabc,bcabca): want 5, got %d", got)
	}
	if got := longestCommonSubstring("", "abc"); got != 0 {
		t.Fatalf("LCS(\"\",abc): want 0, got %d", got)
	}
	// Null-byte regressions: '\x00' in s or t must never collide with internal separator.
	if got := longestCommonSubstring("\x00", "a"); got != 0 {
		t.Fatalf("LCS(\\x00, a): want 0, got %d", got)
	}
	if got := longestCommonSubstring("\x00", "\x00"); got != 1 {
		t.Fatalf("LCS(\\x00, \\x00): want 1, got %d", got)
	}
	if got := longestCommonSubstring("a\x00b", "x\x00by"); got != 2 {
		t.Fatalf("LCS(a\\x00b, x\\x00by): want 2, got %d", got)
	}

	rng := rand.New(rand.NewSource(23))
	bruteLCS := func(s, w string) int {
		best := 0
		for i := 0; i < len(s); i++ {
			for j := 0; j < len(w); j++ {
				k := 0
				for i+k < len(s) && j+k < len(w) && s[i+k] == w[j+k] {
					k++
				}
				if k > best {
					best = k
				}
			}
		}
		return best
	}
	for tc := 0; tc < 60; tc++ {
		b1 := make([]byte, rng.Intn(25))
		b2 := make([]byte, rng.Intn(25))
		for i := range b1 {
			b1[i] = byte(rng.Intn(3)) // includes 0 ('\x00'), 1, 2
		}
		for i := range b2 {
			b2[i] = byte(rng.Intn(3))
		}
		s1, s2 := string(b1), string(b2)
		if got, want := longestCommonSubstring(s1, s2), bruteLCS(s1, s2); got != want {
			t.Fatalf("tc %d LCS(%q,%q): got %d want %d", tc, s1, s2, got, want)
		}
	}
}

func TestBuildSAIntLargeAndNegative(t *testing.T) {
	rng := rand.New(rand.NewSource(31))
	for tc := 0; tc < 50; tc++ {
		n := 1 + rng.Intn(40)
		a := make([]int, n)
		for i := range a {
			a[i] = (rng.Intn(5) - 2) * 1_000_000_000
		}
		sa := buildSAInt(a)
		want := make([]int, n)
		for i := range want {
			want[i] = i
		}
		sort.Slice(want, func(i, j int) bool {
			p, q := want[i], want[j]
			for p < n && q < n && a[p] == a[q] {
				p++
				q++
			}
			if p == n {
				return q != n
			}
			if q == n {
				return false
			}
			return a[p] < a[q]
		})
		for i := range sa {
			if sa[i] != want[i] {
				t.Fatalf("tc %d sa[%d]: got %d want %d", tc, i, sa[i], want[i])
			}
		}
		lcp := buildLCPInt(a, sa)
		for i := 1; i < n; i++ {
			p, q, h := sa[i-1], sa[i], 0
			for p+h < n && q+h < n && a[p+h] == a[q+h] {
				h++
			}
			if lcp[i] != h {
				t.Fatalf("tc %d lcp[%d]: got %d want %d", tc, i, lcp[i], h)
			}
		}
	}
}

func TestBuildSAEdgeAndRandomBruteForce(t *testing.T) {
	if got := buildSA(""); len(got) != 0 {
		t.Fatalf("empty string: want nil/empty, got %v", got)
	}
	if got := buildSA("x"); len(got) != 1 || got[0] != 0 {
		t.Fatalf("single char: want [0], got %v", got)
	}

	rng := rand.New(rand.NewSource(7))
	for tc := 0; tc < 100; tc++ {
		n := 1 + rng.Intn(80)
		alpha := 1 + rng.Intn(5)
		buf := make([]byte, n)
		for i := range buf {
			// Include byte 0 in some tests to ensure separator/null byte safety.
			if tc%5 == 0 {
				buf[i] = byte(rng.Intn(alpha))
			} else {
				buf[i] = byte('a' + rng.Intn(alpha))
			}
		}
		s := string(buf)
		sa := buildSA(s)

		wantSA := make([]int, n)
		for i := range wantSA {
			wantSA[i] = i
		}
		sort.Slice(wantSA, func(i, j int) bool {
			return s[wantSA[i]:] < s[wantSA[j]:]
		})
		for i := 0; i < n; i++ {
			if sa[i] != wantSA[i] {
				t.Fatalf("tc %d sa[%d]: want %d, got %d (s=%q)", tc, i, wantSA[i], sa[i], s)
			}
		}

		lcp := buildLCP(s, sa)
		if lcp[0] != 0 {
			t.Fatalf("tc %d lcp[0]: want 0, got %d", tc, lcp[0])
		}
		for i := 1; i < n; i++ {
			a, b := s[sa[i-1]:], s[sa[i]:]
			want := 0
			for want < len(a) && want < len(b) && a[want] == b[want] {
				want++
			}
			if lcp[i] != want {
				t.Fatalf("tc %d lcp[%d]: want %d, got %d", tc, i, want, lcp[i])
			}
		}
	}
}
