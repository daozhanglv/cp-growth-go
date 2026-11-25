# cp-growth-go

English overview of [cp-growth-go](./README.md): daozhanglv's competitive-programming growth log in Go, with personal solutions and notes plus a set of algorithm templates covered by automated tests.

> This project is based on [codeforces-go](https://github.com/EndlessCheng/codeforces-go) by [EndlessCheng](https://github.com/EndlessCheng). Most templates and historical solutions come from the original project; see [Credits](#credits).

## Quick Start

```bash
# requires the Go version declared in go.mod
git clone https://github.com/daozhanglv/cp-growth-go.git
cd cp-growth-go
sh scripts/check.sh
```

The check script vets and tests the template packages and the shared test utilities. Do not run `go build ./...` or `go test ./...`: the repository holds thousands of single-file solutions that share `package main`, and some tests need the network or local sample files. See [CONTRIBUTING.md](./CONTRIBUTING.md).

---

## Algorithm Index

### Data Structures

| File | Contents |
|------|----------|
| [fenwick_tree.go](copypasta/fenwick_tree.go) | Binary Indexed Tree (BIT / Fenwick Tree), 2-D BIT, offline 2-D point queries |
| [segment_tree.go](copypasta/segment_tree.go) | Segment tree with lazy propagation, persistent (Merge-Sort Tree), dynamic nodes, seg-tree merge/split, offline CDQ |
| [sparse_table.go](copypasta/sparse_table.go) | Sparse Table (RMQ), Disjoint Sparse Table, 2-D Sparse Table |
| [union_find.go](copypasta/union_find.go) | Union-Find with weights, rollback, offline dynamic connectivity |
| [heap.go](copypasta/heap.go) | Augmented heap (modify/delete by key), lazy-deletion heap, regret-greedy heap |
| [treap.go](copypasta/treap.go) | Implicit treap (split/merge), rank tree |
| [splay.go](copypasta/splay.go) | Splay tree (ordered multiset, range reversal split/merge) |
| [skip_list.go](copypasta/skip_list.go) | **[New]** Skip list — O(log n) ordered map, lower/upper bound, floor/ceiling |
| [van_emde_boas.go](copypasta/van_emde_boas.go) | **[New]** Van Emde Boas tree — O(log log U) successor/predecessor |
| [li_chao_tree.go](copypasta/li_chao_tree.go) | **[New]** Li Chao segment tree — dynamic lines, segments, tree merge, max/min queries |
| [wavelet_tree.go](copypasta/wavelet_tree.go) | **[New]** Coordinate-compressed wavelet tree — range k-th smallest, count <= v, count in [lo, hi] |
| [skew_heap.go](copypasta/skew_heap.go) | **[New]** Skew heap — self-adjusting mergeable min-heap |
| [cartesian_tree.go](copypasta/cartesian_tree.go) | Min/Max Cartesian tree construction in O(n) |
| [link_cut_tree.go](copypasta/link_cut_tree.go) | Link-Cut Tree (LCT), dynamic tree connectivity |
| [top_tree.go](copypasta/top_tree.go) | Top Tree |
| [segment_tree_rect.go](copypasta/segment_tree_rect.go) | Sweep-line area of rectangle union |
| [seg_in_bit.go](copypasta/seg_in_bit.go) | Segment tree inside BIT (2-D point update range query) |
| [kd_tree.go](copypasta/kd_tree.go) | k-d Tree, nearest-neighbor search |
| [sqrt_decomposition.go](copypasta/sqrt_decomposition.go) | Sqrt decomposition, block decomposition |
| [mo.go](copypasta/mo.go) | Mo's algorithm (offline, with updates, rollback, on trees) |

### String Algorithms

| File | Contents |
|------|----------|
| [strings.go](copypasta/strings.go) | KMP, Z-algorithm, Manacher, minimal period, string hashing |
| [suffix_array.go](copypasta/suffix_array.go) | **[New]** Suffix array (O(n log n) for strings & int slices), Kasai LCP, distinct substrings, LCS |
| [acam.go](copypasta/acam.go) | Aho-Corasick automaton |
| [sam.go](copypasta/sam.go) | Suffix Automaton (SAM) |
| [pam.go](copypasta/pam.go) | Palindromic Automaton (PAM / Eertree) |
| [trie.go](copypasta/trie.go) | Trie (prefix tree) |
| [trie01.go](copypasta/trie01.go) | Binary trie (XOR trie) |

### Mathematics

| File | Contents |
|------|----------|
| [math.go](copypasta/math.go) | Number theory: primes, factorization, Euler's totient, Möbius function, CRT |
| [math_comb.go](copypasta/math_comb.go) | Combinatorics: binomial coefficients, Catalan, Stirling, inclusion-exclusion |
| [math_matrix.go](copypasta/math_matrix.go) | Matrix exponentiation, Gaussian elimination, Berlekamp-Massey, Kitamasa |
| [berlekamp_massey.go](copypasta/berlekamp_massey.go) | **[New]** Arbitrary-prime Berlekamp-Massey + Kitamasa k-th term over Z/pZ |
| [math_fft.go](copypasta/math_fft.go) | FFT, divide-and-conquer FFT, 2-D FFT |
| [math_ntt.go](copypasta/math_ntt.go) | NTT and polynomial operations (inverse, sqrt, log, exp) |
| [math_gf.go](copypasta/math_gf.go) | Generating functions |
| [math_fwt.go](copypasta/math_fwt.go) | Fast Walsh-Hadamard Transform |
| [math_numerical_analysis.go](copypasta/math_numerical_analysis.go) | Adaptive Simpson's Rule integration, Lagrange polynomial interpolation |
| [math_continued_fraction.go](copypasta/math_continued_fraction.go) | Continued fractions, Pell's equation |
| [big.go](copypasta/big.go) | Big integer and rational arithmetic wrappers |

### Graph Theory

| File | Contents |
|------|----------|
| [graph.go](copypasta/graph.go) | Shortest paths (Dijkstra, Bellman-Ford, SPFA, Floyd), topological sort, SCC, bridges, 2-SAT, bipartite matching (Hopcroft-Karp, Hungarian/KM) |
| [graph_tree.go](copypasta/graph_tree.go) | Tree: LCA, heavy-light decomposition, centroid decomposition, virtual tree, DSU on tree |
| [graph_flow.go](copypasta/graph_flow.go) | Network flow: max flow (Dinic, ISAP, HLPP), min-cost max flow |
| [graph_grid.go](copypasta/graph_grid.go) | Grid graph algorithms |

### Dynamic Programming

| File | Contents |
|------|----------|
| [dp.go](copypasta/dp.go) | Knapsack (0-1, unbounded, grouped), interval DP, digit DP, bitmask DP, tree DP, DP optimizations (monotone queue, divide & conquer, Knuth, slope trick) |
| [convex_hull_trick.go](copypasta/convex_hull_trick.go) | **[New]** Monotone Convex Hull Trick (`chtDeque`, `chtStack`, `chtMinDeque`, `chtMinStack`) |

### Geometry

| File | Contents |
|------|----------|
| [geometry.go](copypasta/geometry.go) | 2-D computational geometry: lines, circles, polygons, convex hull, rotating calipers, half-plane intersection, closest pair |

### Miscellaneous

| File | Contents |
|------|----------|
| [bits.go](copypasta/bits.go) | Bit manipulation tricks |
| [sort.go](copypasta/sort.go) | Binary search, ternary search, 0-1 fractional programming, sorting & inversions |
| [search.go](copypasta/search.go) | Permutations, combinations, Cantor expansion, Gosper's hack, meet-in-the-middle |
| [games.go](copypasta/games.go) | Game theory: Sprague-Grundy, Nim |
| [rand.go](copypasta/rand.go) | Randomized algorithms, simulated annealing |

---

## What This Project Adds

| Addition | Description |
|----------|-------------|
| Template tests | Brute-force comparison tests for union-find, segment trees, Fenwick trees, Sparse Table, shortest paths, modular math, binomials, KMP/Z/Manacher, Aho-Corasick, SAM, and Splay |
| New templates with tests | Skip list, Van Emde Boas tree, Li Chao tree (segments, merge, min/max), coordinate-compressed wavelet tree, skew heap, Berlekamp-Massey + Kitamasa, convex hull trick (min/max), suffix array (strings & ints), Cartesian tree helpers |
| `scripts/check.sh` and CI | One command, also run by GitHub Actions, to format-check, vet, and test templates, utilities, and personal solutions |
| `my-solutions/` | Testable personal Go solutions, notes, and automated index syncing (`scripts/new-solution.sh`, `scripts/sync-solutions.sh`) |
| `docs/` | English documentation site (GitHub Pages) |

---

## Credits

Thanks to [EndlessCheng](https://github.com/EndlessCheng) for open-sourcing [codeforces-go](https://github.com/EndlessCheng/codeforces-go). The templates in `copypasta/` not listed above and the historical solutions in `main/`, `leetcode/` and `misc/` come from the original project. The original copyright notice is kept in [LICENSE](./LICENSE) (MIT).
