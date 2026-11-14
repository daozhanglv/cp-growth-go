---
layout: page
title: Data Structures
permalink: /algorithms/data-structures/
---

# Data Structures

## Skip List

**File:** [`copypasta/skip_list.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/skip_list.go)

A probabilistic sorted-map that supports O(log n) expected operations. Redis uses skip lists for its sorted-set (ZSET) type.

| Operation | Complexity |
|-----------|-----------|
| insert(key, val) | O(log n) expected |
| search(key) | O(log n) expected |
| delete(key) | O(log n) expected |
| lowerBound(x) | O(log n) expected |
| upperBound(x) | O(log n) expected |
| floor(x) | O(log n) expected |

```go
sl := newSkipList()
sl.insert(3, 100)
sl.insert(7, 200)
sl.insert(1, 50)

if v, ok := sl.search(7); ok {
    fmt.Println(v) // 200
}

node := sl.lowerBound(5)
fmt.Println(node.key, node.val) // 7 200

sl.forEach(func(k, v int) {
    fmt.Printf("%d->%d\n", k, v)
})
// 1->50, 3->100, 7->200
```

**Related problems:**
- [LC1206 Design Skiplist](https://leetcode.cn/problems/design-skiplist/)
- [CF 1093G](https://codeforces.com/problemset/problem/1093/G) (ordered multiset, 2300)

---

## Van Emde Boas Tree

**File:** [`copypasta/van_emde_boas.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/van_emde_boas.go)

Maintains a dynamic subset of `[0, U-1]` with O(log log U) per operation — faster than both BSTs and segment trees for predecessor/successor queries on small universes.

| Operation | Complexity |
|-----------|-----------|
| Insert(x) | O(log log U) |
| Delete(x) | O(log log U) |
| Member(x) | O(log log U) |
| Successor(x) | O(log log U) |
| Predecessor(x) | O(log log U) |
| Min() / Max() | O(1) |

Space: O(U) — pre-allocates the full universe. Use when U ≤ 2²⁰.

```go
v := newVEB(1 << 20) // universe [0, 1048575]
v.Insert(42)
v.Insert(17)
v.Insert(100)

fmt.Println(v.Min())          // 17
fmt.Println(v.Successor(42))  // 100
fmt.Println(v.Predecessor(50)) // 42
v.Delete(42)
fmt.Println(v.Member(42))     // false
```

---

## Segment Tree

**File:** [`copypasta/segment_tree.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/segment_tree.go)

Feature-rich segment tree covering: lazy propagation, persistent (merge-sort tree), dynamic nodes, seg-tree merge, segment tree beats, and CDQ divide & conquer.

---

## Binary Indexed Tree (Fenwick Tree)

**File:** [`copypasta/fenwick_tree.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/fenwick_tree.go)

Prefix-sum BIT, range-update BIT, 2-D BIT, offline 2-D point queries.

---

## Union-Find (Disjoint Set Union)

**File:** [`copypasta/union_find.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/union_find.go)

Includes: weighted DSU, rollback DSU, offline dynamic graph connectivity, and bipartite checking.

---

## Convex Hull Trick (CHT)

**File:** [`copypasta/convex_hull_trick.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/convex_hull_trick.go)

Slope-trick DP optimization for transitions of the form `dp[i] = max(k[j]*x[i] + b[j])`.

| Variant | Constraint | Complexity |
|---------|-----------|-----------|
| `chtDeque` | Lines: decreasing k; Queries: increasing x | O(1) amortized |
| `chtStack` | Lines: decreasing k; Queries: arbitrary | O(log n) per query |

Use `chtDeque` when both slopes and query x-values are monotone; otherwise use the Li Chao Tree.

---

## Li Chao Tree

**File:** [`copypasta/li_chao_tree.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/li_chao_tree.go)

Dynamic segment tree for maximum-y queries over lines `y = kx + b`, arbitrary insert and query order.

| Operation | Complexity |
|-----------|-----------|
| AddLine(k, b) | O(log n) |
| QueryMax(x) | O(log n) |

---

## Wavelet Tree

**File:** [`copypasta/wavelet_tree.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/wavelet_tree.go)

Static structure on integer arrays supporting powerful range queries.

| Operation | Complexity |
|-----------|-----------|
| Kth(l, r, k) — k-th smallest in range | O(log U) |
| CountLE(l, r, v) — count ≤ v in range | O(log U) |
| Build | O(n log U) |

---

## Mergeable Heap (Skew Heap)

**File:** [`copypasta/skew_heap.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/skew_heap.go)

Self-adjusting binary min-heap supporting O(log n) amortized merges.

```go
h1 := skewPush(nil, 3)
h1 = skewPush(h1, 1)
h2 := skewPush(nil, 2)
h := skewMerge(h1, h2)  // merged heap, min = 1
fmt.Println(h.val)       // 1
h = skewPop(h)
fmt.Println(h.val)       // 2
```

---

## Cartesian Tree

**File:** [`copypasta/cartesian_tree.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/cartesian_tree.go)

O(n) construction of a min-heap-ordered binary tree that is also BST-ordered by index.
In-order traversal yields the original array indices 0, 1, …, n−1.
The LCA of indices i and j in the Cartesian tree equals the RMQ minimum of a[i..j].
