package ds

type object struct {
	value interface{}
}

func newObject(value interface{}) *object {
	return &object{value: value}
}

type node struct {
	object *object

	next *node
	prev *node
}

func newNode(obj *object, prev, next *node) *node {
	return &node{
		object: obj,
		prev:   prev,
		next:   next,
	}
}

func newHeadNode() *node {
	return &node{nil, nil, nil}
}

// Set represents a data structure of iterable items, grouped
// by key uniqueness
type Set struct {
	items map[object]*node
	head  *node
	foot  *node
}

// NewSet returns an empty fresh set
func NewSet() *Set {
	set := &Set{
		items: make(map[object]*node),
		head:  newHeadNode(),
	}
	set.foot = set.head
	return set
}

func (s *Set) insertHead(n *node) {
	root := s.head
	head := root.next
	if head != nil {
		head.prev = n
	} else {
		s.foot = n
	}
	s.head.next = n
}

// Add inserts into the set the new key-value pair. If the key
// exists already, update the value
func (s *Set) Add(value interface{}) {
	obj := newObject(value)
	if _, ok := s.items[*obj]; ok {
		return
	}
	node := newNode(obj, s.head, s.head.next)
	s.insertHead(node)
	s.items[*obj] = node
}

// Remove erases from the set a pair from a given key. If the
// provided key exists, the pair is returned. Otherwise, nil.
func (s *Set) Remove(value interface{}) *object {
	obj := newObject(value)
	node, ok := s.items[*obj]
	if !ok {
		return nil
	}
	if node == s.foot {
		s.foot = node.prev
	}
	node.prev.next = node.next
	if node.next != nil {
		node.next.prev = node.prev
	}
	object := node.object
	delete(s.items, *object)
	return object
}

// PopFirst removes and returns the first pair that
// entered the set
func (s *Set) PopFirst() interface{} {
	if s.foot == s.head {
		return nil
	}
	obj := s.foot.object
	s.Remove(obj.value)
	return obj.value
}

// Elements returns an iterable slice of de-referenced pairs
func (s *Set) Elements() []object {
	var elements []object
	item := s.head.next
	for item != nil {
		elements = append(elements, *item.object)
		item = item.next
	}
	return elements
}

// Length returns how many elements are in the set
func (s *Set) Length() int {
	return len(s.items)
}
