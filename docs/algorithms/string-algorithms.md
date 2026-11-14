---
layout: page
title: String Algorithms
permalink: /algorithms/string-algorithms/
---

# String Algorithms

## KMP & Z-Algorithm

**File:** [`copypasta/strings.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/strings.go)

| Algorithm | Complexity | Use case |
|-----------|-----------|---------|
| KMP | O(n+m) | Pattern matching, failure function |
| Z-algorithm | O(n) | All occurrences of pattern in text |
| String hashing | O(n) preprocess, O(1) query | Rolling hash, substring comparison |

## Suffix Automaton (SAM)

**File:** [`copypasta/sam.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/sam.go)

Accepts all substrings of a string. Supports:
- Count of distinct substrings — O(n)
- Longest common substring — O(n+m)
- Number of occurrences of each substring
- Lexicographically k-th substring

## Aho-Corasick Automaton

**File:** [`copypasta/acam.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/acam.go)

Multi-pattern matching in O(n + m + k) where n = text length, m = total pattern length, k = number of matches.

## Palindromic Automaton (Eertree)

**File:** [`copypasta/pam.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/pam.go)

Compactly represents all distinct palindromic substrings. Supports:
- Count of palindromic substrings
- Longest palindromic suffix at each position
- Palindrome factorization

## Trie

**File:** [`copypasta/trie.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/trie.go)

Prefix tree for string search / prefix counting.

**File:** [`copypasta/trie01.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/trie01.go)

Binary (XOR) trie — maximum XOR pair, k-th XOR value.
