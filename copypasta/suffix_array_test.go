package copypasta

import (
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
}
