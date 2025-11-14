package treex

// todo 二叉树遍历（三种）

type BinaryTree struct {
	Root *BNode
}

func NewBinaryTree(root *BNode) *BinaryTree {
	return &BinaryTree{Root: root}
}

// PreOrder 二叉树先序遍历, 传入操作节点的函数 operate
func (bt *BinaryTree) PreOrder(operate func(node *BNode)) {
	if bt.Root == nil {
		return
	}
	bt.Root.preOrder(operate)
}

// MiddleOrder 二叉树中序遍历, 传入操作节点的函数 operate
func (bt *BinaryTree) MiddleOrder(operate func(node *BNode)) {
	if bt.Root == nil {
		return
	}
	bt.Root.middleOrder(operate)
}

// PostOrder 二叉树后序遍历, 传入操作节点的函数 operate
func (bt *BinaryTree) PostOrder(operate func(node *BNode)) {
	if bt.Root == nil {
		return
	}
	bt.Root.postOrder(operate)
}

// todo 二叉树层序遍历
