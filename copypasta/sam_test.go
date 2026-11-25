package copypasta

import (
	"math/rand"
	"strings"
	"testing"
)

func TestSAMMatchesBruteForce(t *testing.T) {
	rng := rand.New(rand.NewSource(501))
	for tc := 0; tc < 40; tc++ {
		n := 1 + rng.Intn(30)
		buf := make([]byte, n)
		for i := range buf {
			buf[i] = byte('a' + rng.Intn(4))
		}
		text := string(buf)
		m := newSam(text)

		// 1. Distinct non-empty substrings = sum_{v != root} (v.len - v.fa.len)
		gotDistinct := 0
		for _, o := range m.nodes[1:] {
			gotDistinct += o.len - o.fa.len
		}
		subCount := map[string]int{}
		for i := 0; i < n; i++ {
			for j := i + 1; j <= n; j++ {
				subCount[text[i:j]]++
			}
		}
		if gotDistinct != len(subCount) {
			t.Fatalf("tc %d distinct substrings: got %d, want %d (text=%q)", tc, gotDistinct, len(subCount), text)
		}

		// 2. index(sub) == strings.Index(text, sub) & longestPrefix(q)
		for q := 0; q < 30; q++ {
			qLen := 1 + rng.Intn(10)
			qb := make([]byte, qLen)
			for i := range qb {
				qb[i] = byte('a' + rng.Intn(5)) // includes 'e' not in text
			}
			qs := string(qb)
			if got, want := m.index(qs), strings.Index(text, qs); got != want {
				t.Fatalf("tc %d index(%q) in %q: got %d, want %d", tc, qs, text, got, want)
			}
			wantLP := 0
			for l := 1; l <= len(qs); l++ {
				if strings.Contains(text, qs[:l]) {
					wantLP = l
				}
			}
			if got := m.longestPrefix(qs); got != wantLP {
				t.Fatalf("tc %d longestPrefix(%q) in %q: got %d, want %d", tc, qs, text, got, wantLP)
			}
		}

		// 3. lcs(other) matches longestCommonSubstring
		ob := make([]byte, 1+rng.Intn(25))
		for i := range ob {
			ob[i] = byte('a' + rng.Intn(4))
		}
		other := string(ob)
		if got, want := m.lcs(other), longestCommonSubstring(text, other); got != want {
			t.Fatalf("tc %d sam.lcs(%q) vs %q: got %d, want %d", tc, other, text, got, want)
		}

		// 4. example(): max len(sub)*cnt(sub), tie-break smallest sub (note: example mutates cnt once)
		wantSub := ""
		wantScore := 0
		for sub, cnt := range subCount {
			score := len(sub) * cnt
			if score > wantScore || (score == wantScore && sub < wantSub) {
				wantScore = score
				wantSub = sub
			}
		}
		if got := m.example(); got != wantSub {
			t.Fatalf("tc %d example() on %q: got %q, want %q", tc, text, got, wantSub)
		}
	}
}
