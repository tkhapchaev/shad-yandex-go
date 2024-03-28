//go:build !solution

package lrucache

import "container/list"

type KeyValue struct {
	k int
	v int
}

type LruCache struct {
	capacity int
	data     map[int]*list.Element
	nodes    *list.List
}

func New(capacity int) Cache {
	return &LruCache{
		capacity: capacity,
		data:     make(map[int]*list.Element, capacity),
		nodes:    list.New(),
	}
}

func (lruCache *LruCache) Get(key int) (int, bool) {
	if node, exists := lruCache.data[key]; exists {
		lruCache.nodes.MoveToFront(node)

		return node.Value.(*KeyValue).v, true
	}

	return -1, false
}

func (lruCache *LruCache) Set(key int, value int) {
	if lruCache.capacity == 0 {
		return
	}

	if node, exists := lruCache.data[key]; exists {
		lruCache.nodes.MoveToFront(node)
		node.Value.(*KeyValue).v = value
	} else {
		if lruCache.nodes.Len() == lruCache.capacity {
			index := lruCache.nodes.Back().Value.(*KeyValue).k

			delete(lruCache.data, index)
			lruCache.nodes.Remove(lruCache.nodes.Back())
		}

		newNode := lruCache.nodes.PushFront(&KeyValue{k: key, v: value})
		lruCache.data[key] = newNode
	}
}

func (lruCache *LruCache) Clear() {
	lruCache.data = make(map[int]*list.Element, lruCache.capacity)
	lruCache.nodes = list.New()
}

func (lruCache *LruCache) Range(f func(key, value int) bool) {
	for node := lruCache.nodes.Back(); node != nil; node = node.Prev() {
		if !f(node.Value.(*KeyValue).k, node.Value.(*KeyValue).v) {
			return
		}
	}
}
