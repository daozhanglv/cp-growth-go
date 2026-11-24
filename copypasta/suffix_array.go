package copypasta

/*
Suffix Array (后缀数组)
https://en.wikipedia.org/wiki/Suffix_array
https://oi-wiki.org/string/sa/

See also: strings.go (suffixArray) for an O(n) SA-IS wrapper using stdlib index/suffixarray.

SA[i]  = starting index of the i-th lexicographically smallest suffix
rk[i]  = rank (0-indexed) of suffix starting at i
lcp[i] = length of longest common prefix between SA[i-1] and SA[i]  (lcp[0] = 0)

Build: O(n log n) via prefix doubling + radix (counting) sort
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

// buildSA builds the suffix array of s in O(n log n) using prefix doubling and counting sort.
// Returns sa where sa[i] is the 0-based start of the i-th smallest suffix.
func buildSA(s string) []int {
	n := len(s)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []int{0}
	}
	sa := make([]int, n)
	rk := make([]int, n)
	tmp := make([]int, n)
	m := 256
	if n > m {
		m = n
	}
	cnt := make([]int, m)

	// Initial counting sort by single byte.
	for i := 0; i < n; i++ {
		rk[i] = int(s[i])
		cnt[rk[i]]++
	}
	for i := 1; i < 256; i++ {
		cnt[i] += cnt[i-1]
	}
	for i := n - 1; i >= 0; i-- {
		cnt[rk[i]]--
		sa[cnt[rk[i]]] = i
	}
	tmp[sa[0]] = 0
	classes := 1
	for i := 1; i < n; i++ {
		if rk[sa[i]] != rk[sa[i-1]] {
			classes++
		}
		tmp[sa[i]] = classes - 1
	}
	copy(rk, tmp)

	// Prefix doubling with radix (counting) sort.
	for k := 1; k < n && classes < n; k <<= 1 {
		// tmp holds suffix indices sorted by their second key (i + k).
		p := 0
		for i := n - k; i < n; i++ {
			tmp[p] = i
			p++
		}
		for i := 0; i < n; i++ {
			if sa[i] >= k {
				tmp[p] = sa[i] - k
				p++
			}
		}

		// Stable counting sort by first key rk[tmp[i]].
		clear(cnt[:classes])
		for i := 0; i < n; i++ {
			cnt[rk[tmp[i]]]++
		}
		for i := 1; i < classes; i++ {
			cnt[i] += cnt[i-1]
		}
		for i := n - 1; i >= 0; i-- {
			idx := tmp[i]
			cnt[rk[idx]]--
			sa[cnt[rk[idx]]] = idx
		}

		// Recompute equivalence classes.
		tmp[sa[0]] = 0
		classes = 1
		for i := 1; i < n; i++ {
			cur, prev := sa[i], sa[i-1]
			cur2, prev2 := -1, -1
			if cur+k < n {
				cur2 = rk[cur+k]
			}
			if prev+k < n {
				prev2 = rk[prev+k]
			}
			if rk[cur] != rk[prev] || cur2 != prev2 {
				classes++
			}
			tmp[cur] = classes - 1
		}
		rk, tmp = tmp, rk
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
