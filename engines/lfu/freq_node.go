package engines

type freqNode struct {
	elements map[cacheKey]bool
	// the value representing the frequency
	value uint
	// pointers to compose the dll
	next *freqNode
	prev *freqNode
}

func newHeadFreqNode() *freqNode {
	node := &freqNode{
		value:    0,
		elements: make(map[cacheKey]bool),
	}
	node.prev = nil
	node.next = newFreqNode(1, node, nil)
	return node
}

func newFreqNode(value uint, prev, next *freqNode) *freqNode {
	return &freqNode{
		next:     next,
		prev:     prev,
		value:    value,
		elements: make(map[cacheKey]bool),
	}
}

func (c *freqNode) Add(key cacheKey) {
	c.elements[key] = true
}

func (c *freqNode) Remove(key cacheKey) {
	delete(c.elements, key)
}

func (c *freqNode) Pop() cacheKey {
	if len(c.elements) < 1 {
		return nil
	}
	var res cacheKey
	for key, _ := range c.elements {
		res = key
		delete(c.elements, key)
		break
	}
	return res
}

func (c *freqNode) Size() int {
	return len(c.elements)
}
