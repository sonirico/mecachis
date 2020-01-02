package lfru

import "fmt"

type cacheNode struct {
	key   cacheKey
	value cacheValue

	// pointer to the current frequency node
	parent *freqNode
}

func newCacheNode(key cacheKey, value cacheValue, parent *freqNode) *cacheNode {
	return &cacheNode{key: key, value: value, parent: parent}
}

func (cn *cacheNode) String() string {
	return fmt.Sprintf("<key: %v, value: %v>", cn.key, cn.value)
}
