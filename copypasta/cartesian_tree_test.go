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

func TestCartesianTreeAllVariantsAndMax(t *testing.T) {
	if buildCartesianTreePtr(nil) != nil {
		t.Fatal("empty ptr tree should be nil")
	}
	if nodes, root := buildCartesianTree(nil); nodes != nil || root != -1 {
		t.Fatalf("empty index tree: got (%v, %d)", nodes, root)
	}

	cases := [][]int{
		{3, 2, 6, 1, 9},
		{2, 2, 1, 1, 3, 1},
		{9, 7, 5, 3, 1},
		{1, 3, 2, 3, 1, 4, 2},
	}
	for _, a := range cases {
		ptrRoot := buildCartesianTreePtr(a)
		lr := buildCartesianTree2(a)
		minNodes, minRoot := buildCartesianTree(a)

		// Verify pointer, [][2]int, and cartNode min-trees match identically and satisfy min-heap + inorder.
		var walkPtr func(o *ctNode) []int
		walkPtr = func(o *ctNode) []int {
			if o == nil {
				return nil
			}
			if o.lr[0] != nil && o.val > o.lr[0].val {
				t.Fatalf("min-heap violated on left child: parent=%d child=%d", o.val, o.lr[0].val)
			}
			if o.lr[1] != nil && o.val > o.lr[1].val {
				t.Fatalf("min-heap violated on right child: parent=%d child=%d", o.val, o.lr[1].val)
			}
			lID, rID := -1, -1
			if o.lr[0] != nil {
				lID = o.lr[0].id
			}
			if o.lr[1] != nil {
				rID = o.lr[1].id
			}
			if lr[o.id] != [2]int{lID, rID} || minNodes[o.id].left != lID || minNodes[o.id].right != rID {
				t.Fatalf("structure mismatch at %d: ptr=(%d,%d) lr=%v minNodes=(%d,%d)",
					o.id, lID, rID, lr[o.id], minNodes[o.id].left, minNodes[o.id].right)
			}
			res := walkPtr(o.lr[0])
			res = append(res, o.id)
			res = append(res, walkPtr(o.lr[1])...)
			return res
		}
		if ptrRoot.id != minRoot {
			t.Fatalf("root mismatch: ptr=%d minRoot=%d", ptrRoot.id, minRoot)
		}
		order := walkPtr(ptrRoot)
		for i, id := range order {
			if id != i {
				t.Fatalf("inorder[%d]=%d, want %d", i, id, i)
			}
		}

		// Verify buildCartesianTreeMax satisfies max-heap + inorder.
		maxNodes, maxRoot := buildCartesianTreeMax(a)
		var walkMax func(u int) []int
		walkMax = func(u int) []int {
			if u == -1 {
				return nil
			}
			if l := maxNodes[u].left; l != -1 && maxNodes[u].val < maxNodes[l].val {
				t.Fatalf("max-heap violated on left: parent=%d child=%d", maxNodes[u].val, maxNodes[l].val)
			}
			if r := maxNodes[u].right; r != -1 && maxNodes[u].val < maxNodes[r].val {
				t.Fatalf("max-heap violated on right: parent=%d child=%d", maxNodes[u].val, maxNodes[r].val)
			}
			res := walkMax(maxNodes[u].left)
			res = append(res, u)
			res = append(res, walkMax(maxNodes[u].right)...)
			return res
		}
		maxOrder := walkMax(maxRoot)
		for i, id := range maxOrder {
			if id != i {
				t.Fatalf("max inorder[%d]=%d, want %d", i, id, i)
			}
		}
	}
}
