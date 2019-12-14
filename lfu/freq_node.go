package lfu

type freqNode struct {
	// TODO: Make a ddl of this
	items map[cacheKey]*cacheNode
	// the value representing the frequency
	value uint
	// pointers to compose the dll
	next *freqNode
	prev *freqNode
}

func newHeadFreqNode() *freqNode {
	node := &freqNode{
		value: 0,
		items: make(map[cacheKey]*cacheNode),
	}
	node.prev = nil
	node.next = newFreqNode(1, node, nil)
	return node
}

func newFreqNode(value uint, prev, next *freqNode) *freqNode {
	return &freqNode{
		next:  next,
		prev:  prev,
		value: value,
		items: make(map[cacheKey]*cacheNode),
	}
}

func (c *freqNode) Add(node *cacheNode) {
	c.items[node.key] = node
}

func (c *freqNode) Remove(node *cacheNode) {
	delete(c.items, node.key)
}

// Pop eventually should remove the LRU element from the still
// unimplemented DLL
func (c *freqNode) Pop() *cacheNode {
	if len(c.items) > 0 {
		for _, first := range c.items {
			delete(c.items, first.key)
			return first
		}
	}
	return nil
}

func (c *freqNode) Size() int {
	return len(c.items)
}
