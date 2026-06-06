package lru

import (
	"fmt"
)

const MIN_CACHE_CAPACITY = 4

type listNode struct {
	prev  *listNode
	next  *listNode
	key   string
	value any
}

type Cache struct {
	sHead, sTail listNode
	hash         map[string]*listNode
	capacity     int
	size         int
}

// utility procedures for linkedlist manipulation

func addToHead(node *listNode, sHead *listNode) {
	node.next = sHead.next
	sHead.next.prev = node
	sHead.next = node
	node.prev = sHead
}

func unlink(node *listNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}

// cache procedures

func (c *Cache) evict() {
	listNode := c.sTail.prev
	delete(c.hash, listNode.key)
	unlink(listNode)
	c.size -= 1
}

// public methods
func LRUCache(capacity int) (*Cache, error) { // initialize
	var err error
	cap := capacity
	if capacity < MIN_CACHE_CAPACITY {
		err = fmt.Errorf("Invalid capacity, forcing minimum capacity (%v)", MIN_CACHE_CAPACITY)
		cap = MIN_CACHE_CAPACITY
	}
	m := make(map[string]*listNode, cap)

	init := Cache{capacity: cap, size: 0, hash: m}
	init.sHead.next = &init.sTail
	init.sTail.prev = &init.sHead

	return &init, err
}

func (c *Cache) Put(key string, value any) {
	if oldNode, ok := c.hash[key]; ok {
		oldNode.value = value
		unlink(c.hash[key])
		addToHead(oldNode, &c.sHead)
		return
	}

	node := &listNode{
		value: value,
		key:   key,
	}
	if c.size == c.capacity {
		c.evict()
	}

	c.size += 1
	addToHead(node, &c.sHead)

	c.hash[key] = node
}

func (c *Cache) Get(key string) (any, bool) {
	if node, ok := c.hash[key]; ok {
		unlink(c.hash[key])
		addToHead(node, &c.sHead)
		return node.value, true
	}
	return nil, false
}
