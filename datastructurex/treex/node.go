package treex

// NodeOperationFunc 对二叉树节点的操作函数
type BNodeOperationFunc func(node *BNode)

// BNode 二叉树节点
type BNode struct {
	Val        any
	LeftChind  *BNode
	RightChind *BNode
}

func NewBNode(val any) *BNode {
	return &BNode{
		Val: val,
	}
}

func (node *BNode) preOrder(operate func(node *BNode)) {
	if node == nil {
		return
	}
	operate(node)
	node.LeftChind.preOrder(operate)
	node.RightChind.preOrder(operate)
}

func (node *BNode) middleOrder(operate func(node *BNode)) {
	if node == nil {
		return
	}
	node.LeftChind.middleOrder(operate)
	operate(node)
	node.RightChind.middleOrder(operate)
}

func (node *BNode) postOrder(operate func(node *BNode)) {
	if node == nil {
		return
	}
	node.LeftChind.postOrder(operate)
	node.RightChind.postOrder(operate)
	operate(node)
}
