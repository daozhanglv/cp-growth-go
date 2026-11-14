---
layout: page
title: Graph Theory
permalink: /algorithms/graph-theory/
---

# Graph Theory

## Shortest Paths

**File:** [`copypasta/graph.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/graph.go)

| Algorithm | Complexity | Use case |
|-----------|-----------|---------|
| Dijkstra (heap) | O((V+E) log V) | Non-negative weights |
| Bellman-Ford | O(VE) | Negative weights, detect negative cycles |
| SPFA | O(VE) worst case | Faster in practice on sparse graphs |
| Floyd-Warshall | O(V³) | All-pairs shortest paths |

## Network Flow

**File:** [`copypasta/graph_flow.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/graph_flow.go)

- **Dinic's algorithm** — O(V²E) max flow, O(E√V) on unit graphs
- **Min-cost max flow (MCMF)** — SPFA-based successive shortest paths
- **Bipartite matching** — Hopcroft-Karp O(E√V)
- **Hungarian algorithm** — O(V³) minimum weight bipartite matching

## Tree Algorithms

**File:** [`copypasta/graph_tree.go`](https://github.com/daozhanglv/cp-growth-go/blob/master/copypasta/graph_tree.go)

| Algorithm | Use case |
|-----------|---------|
| LCA (binary lifting) | Lowest common ancestor in O(log n) |
| Heavy-Light Decomposition | Path queries on trees, O(log²n) |
| Centroid Decomposition | Distance queries, path counting |
| DSU on Tree (small-to-large) | Subtree queries |
| Virtual Tree | Queries on a subset of tree nodes |
| Euler Tour + BIT | Subtree sum/update |

## Strongly Connected Components

Kosaraju's and Tarjan's SCC algorithms, bridge finding, articulation points, and 2-SAT are all in `graph.go`.
