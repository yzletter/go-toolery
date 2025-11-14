package listx_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/data_structurex/listx"
)

func TestDoubleList(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}
	l1 := listx.NewDoubleList[int]()
	l2 := listx.NewDoubleListFromSlice(arr)

	l1.InsertToHead(3)
	l1.InsertToHead(2)
	l1.InsertToHead(1)
	l1.InsertToTail(4)
	l1.InsertToTail(5)

	// 打印每个节点
	var printNodeInformation = func(node *listx.ListNode[int]) { fmt.Printf("%v ", node.Val) }

	fmt.Print("正序遍历链表 l1 : ")
	l1.Traverse(printNodeInformation)
	fmt.Println()
	fmt.Print("逆序遍历链表 l2 : ")
	l2.ReverseTraverse(printNodeInformation)
	fmt.Println()
	fmt.Printf("将链表 l1 转为切片: %v\n", l1.Values())
	fmt.Printf("链表 l1 最后一个节点: %v\n", l1.LastNode())
	fmt.Printf("链表 l1 下标为 2 的节点: %v\n", l1.FindNode(2))
	fmt.Printf("链表 l1 下标为 -1 的节点: %v\n", l1.FindNode(-1)) // 预期为空 :

	// 在下标为 3 的节点前后各插一个 6
	node := l1.FindNode(3)
	l1.InsertBefore(6, node)
	l1.InsertAfter(6, node)
	fmt.Println("在下标为 3 的节点前后各插一个 6 : ")
	// 预期结果 : 1 2 3 6 4 6 5
	l1.Traverse(printNodeInformation)
}
