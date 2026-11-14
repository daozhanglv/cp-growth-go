package copypasta

import "sort"

/*
Suffix Array (后缀数组)
https://en.wikipedia.org/wiki/Suffix_array
https://oi-wiki.org/string/sa/

SA[i]  = starting index of the i-th lexicographically smallest suffix
rk[i]  = rank (0-indexed) of suffix starting at i
lcp[i] = length of longest common prefix between SA[i-1] and SA[i]  (lcp[0] = 0)

Build: O(n log n) via prefix doubling + radix sort
LCP:   O(n) via Kasai's algorithm

Key applications:
  - Number of distinct substrings: sum(n - SA[i] - lcp[i]) for i in 0..n-1
  - Longest common substring of two strings (concatenate with separator)
  - Pattern matching: binary search on SA in O(m log n)
  - Repeated substrings, longest repeated substring

Related problems:
https://www.luogu.com.cn/problem/P3809  (SA template)
https://codeforces.com/problemset/problem/235/C  2400
https://codeforces.com/problemset/problem/1063/F  2900
LC1044 https://leetcode.cn/problems/longest-duplicate-substring/
*/

// buildSA builds the suffix array of s in O(n log n).
// Returns sa where sa[i] is the 0-based start of the i-th smallest suffix.
func buildSA(s string) []int {
	n := len(s)
	if n == 0 {
		return nil
	}
	sa := make([]int, n)
	rk := make([]int, n)
	tmp := make([]int, n)
	for i := range sa {
		sa[i] = i
		rk[i] = int(s[i])
	}
	for k := 1; k < n; k <<= 1 {
		rankOf := func(i, shift int) int {
			j := i + shift
			if j >= n {
				return -1
			}
			return rk[j]
		}
		sort.Slice(sa, func(i, j int) bool {
			a, b := sa[i], sa[j]
			if rk[a] != rk[b] {
				return rk[a] < rk[b]
			}
			return rankOf(a, k) < rankOf(b, k)
		})
		tmp[sa[0]] = 0
		for i := 1; i < n; i++ {
			tmp[sa[i]] = tmp[sa[i-1]]
			a, b := sa[i-1], sa[i]
			if rk[a] != rk[b] || rankOf(a, k) != rankOf(b, k) {
				tmp[sa[i]]++
			}
		}
		copy(rk, tmp)
		if rk[sa[n-1]] == n-1 {
			break
		}
	}
	return sa
}

// buildLCP builds the LCP array using Kasai's algorithm in O(n).
// lcp[i] = LCP of sa[i-1] and sa[i]; lcp[0] = 0.
func buildLCP(s string, sa []int) []int {
	n := len(s)
	lcp := make([]int, n)
	rk := make([]int, n)
	for i, v := range sa {
		rk[v] = i
	}
	h := 0
	for i := 0; i < n; i++ {
		if rk[i] > 0 {
			j := sa[rk[i]-1]
			for i+h < n && j+h < n && s[i+h] == s[j+h] {
				h++
			}
			lcp[rk[i]] = h
			if h > 0 {
				h--
			}
		}
	}
	return lcp
}

// countDistinctSubstrings returns the number of distinct non-empty substrings.
func countDistinctSubstrings(s string) int64 {
	n := len(s)
	if n == 0 {
		return 0
	}
	sa := buildSA(s)
	lcp := buildLCP(s, sa)
	var ans int64
	for i, v := range sa {
		ans += int64(n-v) - int64(lcp[i])
	}
	return ans
}

// longestCommonSubstring returns the length of the longest substring common to s and t.
// Uses the "concatenate with separator" trick.
func longestCommonSubstring(s, t string) int {
	sep := string([]byte{0}) // separator not in s or t
	combined := s + sep + t
	sa := buildSA(combined)
	lcp := buildLCP(combined, sa)
	ns := len(s)
	best := 0
	for i := 1; i < len(combined); i++ {
		// sa[i-1] and sa[i] must be from different strings
		a, b := sa[i-1], sa[i]
		aInS := a < ns
		bInS := b < ns
		if aInS != bInS && lcp[i] > best {
			best = lcp[i]
		}
	}
	return best
}
