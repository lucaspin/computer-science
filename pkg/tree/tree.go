package tree

import "fmt"

type TreeNode struct {
	Key    int
	Left   *TreeNode
	Right  *TreeNode
	Parent *TreeNode
}

func CreateTree(keys []int) *TreeNode {
	if len(keys) == 0 {
		return nil
	}

	root := TreeNode{Key: keys[0], Left: nil, Right: nil, Parent: nil}
	keys = keys[1:]

	for _, key := range keys {
		root.Insert(key)
	}

	return &root
}

func (t *TreeNode) Insert(key int) {
	if key <= t.Key {
		if t.Left == nil {
			t.Left = &TreeNode{Key: key, Left: nil, Right: nil, Parent: t}
		} else {
			t.Left.Insert(key)
		}
	} else {
		if t.Right == nil {
			t.Right = &TreeNode{Key: key, Left: nil, Right: nil, Parent: t}
		} else {
			t.Right.Insert(key)
		}
	}
}

func (t *TreeNode) Search(key int) *TreeNode {
	if t.Key == key {
		return t
	}

	if key < t.Key && t.Left != nil {
		return t.Left.Search(key)
	}

	if key > t.Key && t.Right != nil {
		return t.Right.Search(key)
	}

	return nil
}

func (t *TreeNode) Min() *TreeNode {
	node := t
	for node.Left != nil {
		node = node.Left
	}

	return node
}

func (t *TreeNode) Max() *TreeNode {
	node := t
	for node.Right != nil {
		node = node.Right
	}

	return node
}

func (t *TreeNode) IsRoot() bool {
	return t.Parent == nil
}

func (t *TreeNode) IsInLeftSubtree() bool {
	if t.IsRoot() {
		return false
	}

	return t.Parent.Key >= t.Key
}

func (t *TreeNode) IsInRightSubtree() bool {
	if t.IsRoot() {
		return false
	}

	return t.Parent.Key < t.Key
}

func (t *TreeNode) Successor() *TreeNode {
	if t.Right != nil {
		return t.Right.Min()
	}

	// root node with no right subtree has no successor
	if t.IsRoot() {
		return nil
	}

	// node is in a left subtree, successor is its parent
	if t.IsInLeftSubtree() {
		return t.Parent
	}

	// node is in a right subtree
	// we need to find the parent of the lowest ancestor who is not in a right subtree
	ancestor := t
	for ancestor.IsInRightSubtree() {
		ancestor = ancestor.Parent
	}

	return ancestor.Parent
}

func (t *TreeNode) Predecessor() *TreeNode {
	if t.Left != nil {
		return t.Left.Max()
	}

	// root node with no left subtree has no predecessor
	if t.IsRoot() {
		return nil
	}

	// node is in a right subtree, predecessor is its parent
	if t.IsInRightSubtree() {
		return t.Parent
	}

	// node is in a left subtree
	// we need to find the parent of the lowest ancestor who is not in a left subtree
	ancestor := t
	for ancestor.IsInLeftSubtree() {
		ancestor = ancestor.Parent
	}

	return ancestor.Parent
}

func (t *TreeNode) InOrderWalk() {
	if t.Left != nil {
		t.Left.InOrderWalk()
	}

	fmt.Println(t.Key)

	if t.Right != nil {
		t.Right.InOrderWalk()
	}
}

func (t *TreeNode) PreOrderWalk() {
	fmt.Println(t.Key)

	if t.Left != nil {
		t.Left.PreOrderWalk()
	}

	if t.Right != nil {
		t.Right.PreOrderWalk()
	}
}

func (t *TreeNode) PostOrderWalk() {
	if t.Left != nil {
		t.Left.PostOrderWalk()
	}

	if t.Right != nil {
		t.Right.PostOrderWalk()
	}

	fmt.Println(t.Key)
}
