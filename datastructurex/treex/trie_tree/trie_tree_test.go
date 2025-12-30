package trie_tree_test

import (
	"fmt"
	"testing"

	"github.com/yzletter/go-toolery/datastructurex/treex/trie_tree"
)

func TestTrieTree(t *testing.T) {
	tree := trie_tree.NewTrieTree()

	tree.Add("分布式")
	tree.Add("Golang")
	tree.Add("Go")
	tree.Add("分布式搜索")
	tree.Add("分布式搜索引擎")

	fmt.Println(tree.Retrieve("G"))
	fmt.Println(tree.Retrieve("分布式搜索"))
	fmt.Println(tree.Retrieve("分"))
}

// go test -v ./datastructurex/treex/trie_tree -run=^TestTrieTree$ -count=1
/*
yzletter@yangzhileideMacBook-Pro go-toolery % go test -v ./datastructurex/treex/trie_tree -run=^TestTrieTree$ -count=1
=== RUN   TestTrieTree
[Go Golang]
[分布式搜索 分布式搜索引擎]
[分布式 分布式搜索 分布式搜索引擎]
--- PASS: TestTrieTree (0.00s)
PASS
ok      github.com/yzletter/go-toolery/datastructurex/treex/trie_tree   0.569s
*/
