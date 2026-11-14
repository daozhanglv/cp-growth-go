---
layout: home
title: cp-growth-go
---

# cp-growth-go

daozhanglv's competitive-programming growth log in Go: personal solutions and notes, plus algorithm templates covered by automated tests.

This project is based on [codeforces-go](https://github.com/EndlessCheng/codeforces-go) by EndlessCheng. Most templates and historical solutions come from the original project.

## What is here

- Algorithm templates in `copypasta/`, with reference links and problem lists from the original project
- Automated tests for a growing subset of the templates, compared against brute-force implementations
- New templates with tests: skip list, Van Emde Boas tree, Li Chao tree, wavelet tree, skew heap, Berlekamp-Massey, convex hull trick, suffix array
- A single check command, also run in CI

## Algorithm Categories

| Category | Key Templates |
|----------|--------------|
| [Data Structures](algorithms/data-structures) | Segment tree, BIT, Skip list, VEB tree, LCT, Treap |
| [Graph Theory](algorithms/graph-theory) | Dijkstra, Max Flow, LCA, HLD, SCC, 2-SAT |
| [Dynamic Programming](algorithms/dynamic-programming) | Knapsack, Slope Trick, Digit DP, Bitmask DP |
| [String Algorithms](algorithms/string-algorithms) | KMP, SAM, Aho-Corasick, Suffix Array |
| Mathematics | FFT/NTT, Number Theory, Matrix Exponentiation |

## Quick Start

```bash
git clone https://github.com/daozhanglv/cp-growth-go.git
cd cp-growth-go
sh scripts/check.sh
```

## License

[MIT](https://github.com/daozhanglv/cp-growth-go/blob/master/LICENSE). The original copyright notice by EndlessCheng is kept.
