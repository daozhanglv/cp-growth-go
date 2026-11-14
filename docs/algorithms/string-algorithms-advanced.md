---
layout: page
title: Advanced String & Sequence Algorithms
permalink: /algorithms/string-algorithms-advanced/
---

# Advanced String & Sequence Algorithms

## Suffix Array

**File:** [`copypasta/suffix_array.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/suffix_array.go)

| Function | Description | Complexity |
|----------|-------------|-----------|
| `buildSA(s)` | Build suffix array | O(n log n) |
| `buildLCP(s, sa)` | Build LCP array (Kasai) | O(n) |
| `countDistinctSubstrings(s)` | Count distinct non-empty substrings | O(n log n) |
| `longestCommonSubstring(s, t)` | Longest common substring | O((n+m) log(n+m)) |

`SA[i]` is the start index of the i-th lexicographically smallest suffix.
`LCP[i]` is the length of the longest common prefix between `SA[i-1]` and `SA[i]`.

**Related problems:**
- [Luogu P3809](https://www.luogu.com.cn/problem/P3809) — SA template
- [LC1044](https://leetcode.cn/problems/longest-duplicate-substring/) — Longest Duplicate Substring
- [CF 235C](https://codeforces.com/problemset/problem/235/C) — 2400

## Berlekamp-Massey (Linear Recurrence Finding)

**File:** [`copypasta/berlekamp_massey.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/berlekamp_massey.go)

Given a sequence `s[0..n-1]` over Z/pZ, finds the shortest linear recurrence automatically.

```go
mod := int64(1e9 + 7)
// Fibonacci: 0,1,1,2,3,5,...
s := []int64{0, 1, 1, 2, 3, 5, 8, 13}
rec := berlekampMassey(s, mod)
// rec = [1, 1]  meaning s[i] = 1*s[i-1] + 1*s[i-2]
```

Combine with matrix exponentiation to evaluate `s[n]` for n up to 10^18.

**Related problems:**
- [CF 1511G](https://codeforces.com/problemset/problem/1511/G) — 3100
- [Luogu P5487](https://www.luogu.com.cn/problem/P5487) — BM + Kitamasa
