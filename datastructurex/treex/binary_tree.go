package treex

import (
	"github.com/yzletter/go-toolery/datastructurex/dequeuex"
)

// BinaryTree 二叉树
type BinaryTree struct {
	Root *BNode
}

// NewBinaryTree 根据 root 构造一颗二叉树
func NewBinaryTree(root *BNode) *BinaryTree {
	return &BinaryTree{Root: root}
}

// PreOrder 二叉树先序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) PreOrder(operate BNodeOperationFunc) {
	if bt.Root == nil {
		return
	}
	bt.Root.preOrder(operate)
}

// MiddleOrder 二叉树中序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) MiddleOrder(operate BNodeOperationFunc) {
	if bt.Root == nil {
		return
	}
	bt.Root.middleOrder(operate)
}

// PostOrder 二叉树后序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) PostOrder(operate BNodeOperationFunc) {
	if bt.Root == nil {
		return
	}
	bt.Root.postOrder(operate)
}

// LevelOrder 二叉树层序遍历, 传入操作节点的函数 BNodeOperationFunc
func (bt *BinaryTree) LevelOrder(operate BNodeOperationFunc) {
	q := dequeuex.NewDequeue[*BNode]()
	q.PushBack(bt.Root)
	for q.Size() > 0 {
		node, err := q.Front()
		if err != nil {
			return
		}
		err = q.PopFront()
		if err != nil {
			return
		}

		operate(node)
		if node.LeftChind != nil {
			q.PushBack(node.LeftChind)
		}
		if node.RightChind != nil {
			q.PushBack(node.RightChind)
		}
	}
}
