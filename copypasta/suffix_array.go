package copypasta

import "sort"

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
  - Longest common substring of two strings (concatenate with sentinel)
  - Pattern matching: binary search on SA in O(m log n)
  - Integer-array suffix queries via buildSAInt / buildLCPInt

Related problems:
https://www.luogu.com.cn/problem/P3809  (SA template)
https://codeforces.com/problemset/problem/235/C  2400
https://codeforces.com/problemset/problem/1063/F  2900
LC1044 https://leetcode.cn/problems/longest-duplicate-substring/
*/

// buildSANonNeg builds the suffix array for non-negative integers in [0, alphaSize).
func buildSANonNeg(a []int, alphaSize int) []int {
	n := len(a)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []int{0}
	}
	sa := make([]int, n)
	rk := make([]int, n)
	tmp := make([]int, n)
	m := alphaSize
	if n > m {
		m = n
	}
	cnt := make([]int, m)

	for i, v := range a {
		rk[i] = v
		cnt[v]++
	}
	for i := 1; i < alphaSize; i++ {
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
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = int(s[i])
	}
	return buildSANonNeg(a, 256)
}

// buildSAInt builds the suffix array of an arbitrary integer slice (supports negative and large values).
func buildSAInt(a []int) []int {
	n := len(a)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return []int{0}
	}
	minV, maxV := a[0], a[0]
	for _, v := range a[1:] {
		if v < minV {
			minV = v
		}
		if v > maxV {
			maxV = v
		}
	}
	if minV >= 0 && maxV <= max(256, 2*n) {
		return buildSANonNeg(a, maxV+1)
	}
	vals := append([]int(nil), a...)
	sort.Ints(vals)
	uniq := vals[:0]
	for _, v := range vals {
		if len(uniq) == 0 || uniq[len(uniq)-1] != v {
			uniq = append(uniq, v)
		}
	}
	mapped := make([]int, n)
	for i, v := range a {
		mapped[i] = sort.SearchInts(uniq, v)
	}
	return buildSANonNeg(mapped, len(uniq))
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

// buildLCPInt builds the LCP array for an integer slice and its suffix array in O(n).
func buildLCPInt(a []int, sa []int) []int {
	n := len(a)
	lcp := make([]int, n)
	rk := make([]int, n)
	for i, v := range sa {
		rk[v] = i
	}
	h := 0
	for i := 0; i < n; i++ {
		if rk[i] > 0 {
			j := sa[rk[i]-1]
			for i+h < n && j+h < n && a[i+h] == a[j+h] {
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
// Uses an integer sentinel 0 while shifting all byte values to [1, 256], so null bytes
// ('\x00') inside s or t are handled safely without colliding with the separator.
func longestCommonSubstring(s, t string) int {
	if len(s) == 0 || len(t) == 0 {
		return 0
	}
	ns := len(s)
	combined := make([]int, ns+1+len(t))
	for i := 0; i < ns; i++ {
		combined[i] = int(s[i]) + 1
	}
	combined[ns] = 0 // unique sentinel strictly smaller than all mapped bytes
	for i := 0; i < len(t); i++ {
		combined[ns+1+i] = int(t[i]) + 1
	}
	sa := buildSANonNeg(combined, 257)
	lcp := buildLCPInt(combined, sa)
	best := 0
	for i := 1; i < len(combined); i++ {
		a, b := sa[i-1], sa[i]
		if a == ns || b == ns {
			continue
		}
		if (a < ns) != (b < ns) && lcp[i] > best {
			best = lcp[i]
		}
	}
	return best
}
