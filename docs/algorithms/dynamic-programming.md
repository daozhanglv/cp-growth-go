---
layout: page
title: Dynamic Programming
permalink: /algorithms/dynamic-programming/
---

# Dynamic Programming

**File:** [`copypasta/dp.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/dp.go)

## Knapsack

| Variant | Description |
|---------|-------------|
| 0-1 Knapsack | Each item used at most once — O(nW) |
| Unbounded Knapsack | Items can be reused — O(nW) |
| Grouped Knapsack | Items in groups, pick at most one per group |
| Tree Knapsack | Items form a tree; subtree constraints |

## DP Optimizations

| Technique | Reduces | Condition |
|-----------|---------|-----------|
| Monotone Queue Optimization | O(n²) → O(n) | Transition depends on a sliding window |
| Divide & Conquer Optimization | O(n²) → O(n log n) | Optimal split point is monotone |
| Knuth's Optimization | O(n³) → O(n²) | Interval DP with monotone opt |
| Slope Trick | O(n²) → O(n log n) | DP on piecewise-linear convex functions |
| Convex Hull Trick (CHT) | O(n²) → O(n) | Linear transition cost, monotone slope |

## Other Patterns

- **Digit DP** — count integers in [L, R] satisfying digit constraints
- **Bitmask DP** — subsets of a small set (n ≤ 20)
- **Interval DP** — optimal partitioning of sequences/strings
- **Profile DP** — tiling problems with broken profile
- **Tree DP** — subtree states; rerooting technique
