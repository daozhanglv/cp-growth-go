package copypasta

/*
Cartesian Tree (笛卡尔树)
https://en.wikipedia.org/wiki/Cartesian_tree
https://oi-wiki.org/ds/cartesian-tree/

A binary tree on array a[0..n-1] satisfying:
  - BST by index: in-order traversal yields 0,1,...,n-1
  - Min-heap by value: parent.val ≤ children.val

Construction: O(n) using a monotone stack.

Applications:
  - Range Minimum Query (RMQ) ↔ LCA on the Cartesian tree
  - Treap building block
  - DP on subarray structures

Related problems:
https://www.luogu.com.cn/problem/P5854  (Cartesian tree template)
https://codeforces.com/problemset/problem/1017/E  2500
LC1944 https://leetcode.cn/problems/number-of-visible-people-in-a-queue/

todo 一些题目 https://www.luogu.com.cn/blog/AAAbbb123/di-ka-er-shu-xue-xi-bi-ji
另见单调栈 monotone_stack.go

模板题 https://www.luogu.com.cn/problem/P5854 https://judge.yosupo.jp/problem/cartesian_tree
      https://www.luogu.com.cn/problem/P1377
https://codeforces.com/problemset/problem/1777/D 1900
https://codeforces.com/problemset/problem/1580/B 2600 思想
*/

// ctNode is the pointer-based Cartesian tree node (max-heap by default based on construction).
type ctNode struct {
	lr      [2]*ctNode
	id, val int
}

func buildCartesianTreePtr(a []int) *ctNode {
	if len(a) == 0 {
		return nil
	}
	s := []*ctNode{}
	for i, v := range a {
		o := &ctNode{id: i, val: v}
		for len(s) > 0 {
			top := s[len(s)-1]
			if top.val < v {
				top.lr[1] = o
				break
			}
			o.lr[0] = top
			s = s[:len(s)-1]
		}
		s = append(s, o)
	}
	return s[0]
}

// buildCartesianTree2 returns each node's left/right child indices (-1 if absent).
// 非指针版，返回每个节点的左右儿子的编号
func buildCartesianTree2(a []int) [][2]int {
	n := len(a)
	lr := make([][2]int, n)
	for i := range lr {
		lr[i] = [2]int{-1, -1}
	}
	s := []int{}
	for i, v := range a {
		for len(s) > 0 {
			topI := s[len(s)-1]
			if a[topI] < v {
				lr[topI][1] = i
				break
			}
			lr[i][0] = topI
			s = s[:len(s)-1]
		}
		s = append(s, i)
	}
	return lr
}

// cartNode is used by buildCartesianTree (index-based min-heap Cartesian tree).
type cartNode struct {
	val         int
	left, right int // indices into the nodes slice; -1 = absent
}

// buildCartesianTree builds a min-Cartesian tree from a[0..n-1] in O(n).
// Returns (nodes, root) where root is the index of the root node.
// nodes[i].left and nodes[i].right give children indices (-1 if absent).
func buildCartesianTree(a []int) (nodes []cartNode, root int) {
	n := len(a)
	if n == 0 {
		return nil, -1
	}
	nodes = make([]cartNode, n)
	for i := range nodes {
		nodes[i] = cartNode{val: a[i], left: -1, right: -1}
	}
	stack := make([]int, 0, n)
	for i := 0; i < n; i++ {
		last := -1
		for len(stack) > 0 && nodes[stack[len(stack)-1]].val >= nodes[i].val {
			last = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
		}
		nodes[i].left = last
		if len(stack) > 0 {
			nodes[stack[len(stack)-1]].right = i
		}
		stack = append(stack, i)
	}
	return nodes, stack[0]
}
