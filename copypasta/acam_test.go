package copypasta

import (
	"math/rand"
	"slices"
	"strings"
	"testing"
)

func TestACAMSearchAndSumCount(t *testing.T) {
	patterns := []string{"abc", "abac", "babc", "bc", "cab"}
	text := "abcabacbabccab"
	ac := newACAM(patterns)
	gotPos := ac.acSearch(text)
	for i, p := range patterns {
		var want []int
		for j := 0; j+len(p) <= len(text); j++ {
			if text[j:j+len(p)] == p {
				want = append(want, j)
			}
		}
		if !slices.Equal(gotPos[i], want) {
			t.Fatalf("pattern %q: got %v, want %v", p, gotPos[i], want)
		}
	}

	// Duplicate patterns in sumCountAllPatterns
	dupPatterns := []string{"a", "aa", "aa", "abbc", "babc"}
	text2 := "abcbabcbabcabaacbababacbabc"
	ac2 := newACAM(dupPatterns)
	wantCnt := 0
	for _, p := range dupPatterns {
		if strings.Contains(text2, p) {
			wantCnt++
		}
	}
	if got := ac2.sumCountAllPatterns(text2); got != wantCnt {
		t.Fatalf("sumCountAllPatterns: got %d, want %d", got, wantCnt)
	}
}

func TestACAMRandomizedBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(701))
	for tc := 0; tc < 40; tc++ {
		seen := map[string]bool{}
		var patterns []string
		for len(patterns) < 12 {
			b := make([]byte, 1+rng.Intn(6))
			for i := range b {
				b[i] = byte('a' + rng.Intn(4))
			}
			s := string(b)
			if !seen[s] {
				seen[s] = true
				patterns = append(patterns, s)
			}
		}
		tb := make([]byte, 40+rng.Intn(40))
		for i := range tb {
			tb[i] = byte('a' + rng.Intn(4))
		}
		text := string(tb)

		ac := newACAM(patterns)
		gotPos := ac.acSearch(text)
		wantOccurred := 0
		for i, p := range patterns {
			var want []int
			for j := 0; j+len(p) <= len(text); j++ {
				if text[j:j+len(p)] == p {
					want = append(want, j)
				}
			}
			if len(want) > 0 {
				wantOccurred++
			}
			if !slices.Equal(gotPos[i], want) {
				t.Fatalf("tc %d pattern %q: got %v, want %v", tc, p, gotPos[i], want)
			}
		}
		if got := ac.sumCountAllPatterns(text); got != wantOccurred {
			t.Fatalf("tc %d sumCountAllPatterns: got %d, want %d", tc, got, wantOccurred)
		}
	}
}
