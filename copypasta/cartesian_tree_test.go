package copypasta

import "testing"

func TestCartesianTreeBasic(t *testing.T) {
	// a = [3, 2, 6, 1, 9]
	// min-Cartesian tree: root = 3 (value 1, index 3)
	//     3(1)
	//    /    \
	//  1(2)   4(9)
	//  /   \
	// 0(3) 2(6)
	a := []int{3, 2, 6, 1, 9}
	nodes, root := buildCartesianTree(a)

	if nodes[root].val != 1 {
		t.Fatalf("root val: want 1, got %d", nodes[root].val)
	}
	if root != 3 {
		t.Fatalf("root index: want 3, got %d", root)
	}
	// root's left child should be index 1 (val 2)
	if nodes[root].left != 1 {
		t.Fatalf("root.left: want 1, got %d", nodes[root].left)
	}
	// root's right child should be index 4 (val 9)
	if nodes[root].right != 4 {
		t.Fatalf("root.right: want 4, got %d", nodes[root].right)
	}
}

func TestCartesianTreeSorted(t *testing.T) {
	// sorted ascending → root is index 0 (smallest), right-spine only
	a := []int{1, 2, 3, 4, 5}
	nodes, root := buildCartesianTree(a)
	if root != 0 {
		t.Fatalf("sorted: root should be 0, got %d", root)
	}
	// each node's left should be -1; right should be next
	cur := root
	for i := 0; i < len(a)-1; i++ {
		if nodes[cur].left != -1 {
			t.Fatalf("node %d should have no left child", cur)
		}
		if nodes[cur].right != i+1 {
			t.Fatalf("node %d right: want %d, got %d", cur, i+1, nodes[cur].right)
		}
		cur = nodes[cur].right
	}
}

func TestCartesianTreeInorder(t *testing.T) {
	// in-order traversal must yield indices 0,1,...,n-1
	a := []int{5, 3, 7, 1, 4, 2, 6}
	nodes, root := buildCartesianTree(a)

	var inorder func(i int) []int
	inorder = func(i int) []int {
		if i == -1 {
			return nil
		}
		res := inorder(nodes[i].left)
		res = append(res, i)
		res = append(res, inorder(nodes[i].right)...)
		return res
	}
	got := inorder(root)
	for idx, v := range got {
		if v != idx {
			t.Fatalf("inorder[%d] = %d, want %d", idx, v, idx)
		}
	}
}
