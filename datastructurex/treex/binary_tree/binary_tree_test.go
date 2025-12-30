package binary_tree_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/datastructurex/treex/binary_tree"
)

func TestBinaryTree(t *testing.T) {
	node := &binary_tree.BNode{
		Val: 5,
	}
	root := binary_tree.NewBinaryTree(node)
	node3 := binary_tree.NewBNode(3)
	node4 := binary_tree.NewBNode(4)
	node1 := binary_tree.NewBNode(1)
	node6 := binary_tree.NewBNode(6)
	node7 := binary_tree.NewBNode(7)
	node9 := binary_tree.NewBNode(9)

	root.Root.LeftChind = node3
	root.Root.RightChind = node4
	node3.LeftChind = node1
	node3.RightChind = node6
	node4.LeftChind = node7
	node4.RightChind = node9

	//				5
	//		3				4
	//	1		6		7		9
	var printNode = func(node *binary_tree.BNode) {
		fmt.Printf("%v ", node.Val)
	}
	root.PreOrder(printNode) // 前序 : 5 3 1 6 4 7 9
	fmt.Println()
	root.MiddleOrder(printNode) // 中序 : 1 3 6 5 7 4 9
	fmt.Println()
	root.PostOrder(printNode) // 后序 : 1 6 3 7 9 4 5
	fmt.Println()
	root.LevelOrder(printNode)
}

// go test -v ./data_structure/treex -run=^TestBinaryTree$ -count=1
