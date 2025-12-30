package trie_tree

type TrieTree struct {
	root *TrieNode
}

func NewTrieTree() *TrieTree {
	return &TrieTree{
		root: new(TrieNode), // 虚根
	}
}

// Add 将字符串添加到 Trie 树
func (tree *TrieTree) Add(term string) {
	if tree.root == nil {
		tree.root = new(TrieNode)
	}

	// 必须是词，不能是单个字符
	if len(term) <= 1 {
		return
	}

	tree.root.add(term, 0)
}

// Retrieve 查询符合前缀的所有字符串
func (tree *TrieTree) Retrieve(prefix string) []string {
	if tree.root == nil || len(tree.root.Children) == 0 {
		return nil
	}

	firstWord := []rune(prefix)[0]

	if childNode, exists := tree.root.Children[firstWord]; exists {
		// 查找 prefix 的终点
		prefixEndNode := childNode.walk(prefix, 0)
		if prefixEndNode == nil {
			// 不存在以 prefix 为前缀的字符串
			return nil
		}

		// 返回该节点子树的所有 term
		res := make([]string, 0, 100)
		prefixEndNode.traverseTerms(&res)
		return res
	}
	return nil
}
