package trie_tree

type TrieNode struct {
	Word     rune
	Children map[rune]*TrieNode
	Term     string
}

func (node *TrieNode) add(term string, idx int) {
	if idx >= len([]rune(term)) {
		node.Term = term
		return
	}

	// 没有子节点
	if node.Children == nil {
		node.Children = make(map[rune]*TrieNode)
	}

	word := []rune(term)[idx]
	if childNode, exist := node.Children[word]; !exist {
		// 子节点不存在
		newChild := &TrieNode{Word: word}
		node.Children[word] = newChild
		newChild.add(term, idx+1)
	} else {
		// 子节点存在
		childNode.add(term, idx+1)
	}
}

func (node *TrieNode) walk(prefix string, idx int) *TrieNode {
	// 当前字符为 []rune(prefix)[idx]
	if idx == len([]rune(prefix))-1 {
		return node
	}

	// 下一个字符
	idx++
	word := []rune(prefix)[idx]

	if childNode, exist := node.Children[word]; exist {
		return childNode.walk(prefix, idx)
	}
	return nil
}

func (node *TrieNode) traverseTerms(res *[]string) {
	if len(node.Term) > 0 {
		*res = append(*res, node.Term)
	}

	for _, childNode := range node.Children {
		childNode.traverseTerms(res)
	}
}
