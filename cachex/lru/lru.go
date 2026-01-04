package lru

import (
	"cmp"

	"github.com/yzletter/go-toolery/datastructurex/listx"
)

type LRUCache[K cmp.Ordered, V any] struct {
	cap int
	mp  map[K]V
	l   *listx.LinkedList[K]
}

func NewLRUCache[K cmp.Ordered, V any](cap int) *LRUCache[K, V] {
	mp := make(map[K]V, cap)
	l := listx.NewLinkedList[K]()
	return &LRUCache[K, V]{
		cap: cap,
		mp:  mp,
		l:   l,
	}
}

func (cache *LRUCache[K, V]) Add(key K, value V) {
	if len(cache.mp) == cache.cap {
		node, _ := cache.l.LastNode()
		delete(cache.mp, node.Val)
		_ = cache.l.DeleteLastNode()
	}
	cache.mp[key] = value
	cache.l.InsertToHead(key)
}

func (cache *LRUCache[K, V]) Get(key K) (V, bool) {
	value, ok := cache.mp[key]
	if !ok {
		// 缓存中不存在
		return value, false
	}

	node, _ := cache.l.FindNodeByValue(key)
	cache.l.MoveToHead(node)
	return value, ok
}

func (cache *LRUCache[K, V]) Size() int {
	return len(cache.mp)
}
