package copypasta

import "math/rand"

/*
Skip List 跳表
https://en.wikipedia.org/wiki/Skip_list
Original paper: https://dl.acm.org/doi/10.1145/78973.78977

A probabilistic sorted-map that supports O(log n) expected-time operations.
Each element is stored at one or more levels; higher levels act as express lanes.

Advantages over balanced BSTs:
  - Simpler implementation (no rotations)
  - Redis uses skip lists for ZSET (sorted sets) due to easy range scans

Time complexity (expected): O(log n) for insert / search / delete / range queries
Space complexity (expected): O(n log n)

Go's standard library has no ordered map; this template fills that gap.

Related problems:
LC1206 https://leetcode.cn/problems/design-skiplist/
https://codeforces.com/problemset/problem/61/E
https://codeforces.com/problemset/problem/1093/G  2300  (ordered multiset operations)
*/

const (
	_slMaxLevel = 32   // max level — supports ~4 billion elements
	_slP        = 0.25 // probability that a node promotes to the next level
)

type slNode struct {
	key     int
	val     int
	forward []*slNode
}

// skipList is an ordered map (int key → int value).
// Inserting a duplicate key updates the value in-place.
type skipList struct {
	head  *slNode
	level int // number of active levels (1-indexed)
	size  int
}

func newSkipList() *skipList {
	return &skipList{
		head:  &slNode{forward: make([]*slNode, _slMaxLevel)},
		level: 1,
	}
}

func (*skipList) randomLevel() int {
	lv := 1
	for lv < _slMaxLevel && rand.Float64() < _slP {
		lv++
	}
	return lv
}

// insert inserts key→val; if key exists, the value is updated.
func (sl *skipList) insert(key, val int) {
	var update [_slMaxLevel]*slNode
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && cur.forward[i].key < key {
			cur = cur.forward[i]
		}
		update[i] = cur
	}
	if nxt := update[0].forward[0]; nxt != nil && nxt.key == key {
		nxt.val = val // update existing
		return
	}
	lv := sl.randomLevel()
	if lv > sl.level {
		for i := sl.level; i < lv; i++ {
			update[i] = sl.head
		}
		sl.level = lv
	}
	node := &slNode{key: key, val: val, forward: make([]*slNode, lv)}
	for i := range lv {
		node.forward[i] = update[i].forward[i]
		update[i].forward[i] = node
	}
	sl.size++
}

// search returns (val, true) if key is found, else (0, false).
func (sl *skipList) search(key int) (int, bool) {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && cur.forward[i].key < key {
			cur = cur.forward[i]
		}
	}
	if nxt := cur.forward[0]; nxt != nil && nxt.key == key {
		return nxt.val, true
	}
	return 0, false
}

// delete removes key. Returns true if the key existed.
func (sl *skipList) delete(key int) bool {
	var update [_slMaxLevel]*slNode
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && cur.forward[i].key < key {
			cur = cur.forward[i]
		}
		update[i] = cur
	}
	target := update[0].forward[0]
	if target == nil || target.key != key {
		return false
	}
	for i := range sl.level {
		if update[i].forward[i] != target {
			break
		}
		update[i].forward[i] = target.forward[i]
	}
	for sl.level > 1 && sl.head.forward[sl.level-1] == nil {
		sl.level--
	}
	sl.size--
	return true
}

// lowerBound returns the first node with key >= target, or nil if none.
func (sl *skipList) lowerBound(target int) *slNode {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && cur.forward[i].key < target {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// upperBound returns the first node with key > target, or nil if none.
func (sl *skipList) upperBound(target int) *slNode {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && cur.forward[i].key <= target {
			cur = cur.forward[i]
		}
	}
	return cur.forward[0]
}

// floor returns the last node with key <= target, or nil if none.
func (sl *skipList) floor(target int) *slNode {
	cur := sl.head
	for i := sl.level - 1; i >= 0; i-- {
		for cur.forward[i] != nil && cur.forward[i].key <= target {
			cur = cur.forward[i]
		}
	}
	if cur == sl.head {
		return nil
	}
	return cur
}

// ceiling is an alias for lowerBound (first key >= target).
func (sl *skipList) ceiling(target int) *slNode { return sl.lowerBound(target) }

// forEach calls f(key, val) on every element in ascending order.
func (sl *skipList) forEach(f func(key, val int)) {
	for cur := sl.head.forward[0]; cur != nil; cur = cur.forward[0] {
		f(cur.key, cur.val)
	}
}

func (sl *skipList) len() int { return sl.size }
