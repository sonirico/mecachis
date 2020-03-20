package mecachis

type Cacheable interface {
	Len() int
}

type Value interface {
	Cacheable

	Value() interface{}
}

type Entry interface {
	Key() string
	Value() Value
	Len() uint64
}

type EvictionFn func(Entry)

type Engine interface {
	Insert(k string, v Cacheable) bool
	Upsert(k string, v Cacheable)
	Access(k string, v Cacheable) (Value, bool)
	Size() uint64
	Dump() []Entry
	OnEvict(func(Entry))
}
