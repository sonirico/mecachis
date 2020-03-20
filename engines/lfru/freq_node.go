package engines

import "github.com/sonirico/mecachis/container"

type freqNode struct {
	set *container.Set
	// the value representing the frequency
	value uint
	// pointers to compose the dll
	next *freqNode
	prev *freqNode
}

func newHeadFreqNode() *freqNode {
	node := &freqNode{
		value: 0,
		set:   container.NewSet(),
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
		set:   container.NewSet(),
	}
}

func (c *freqNode) Add(node *cacheNode) {
	c.set.Add(node)
}

func (c *freqNode) Remove(node *cacheNode) {
	c.set.Remove(node.key)
}

func (c *freqNode) Pop() *cacheNode {
	if c.set.Length() < 1 {
		return nil
	}

	lru := c.set.PopFirst()
	node, _ := lru.(*cacheNode)
	return node
}

func (c *freqNode) Size() int {
	return c.set.Length()
}
