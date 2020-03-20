package engines

type Cacheable interface {
	Len() uint64
}

type Value interface {
	Cacheable

	Value() interface{}
}

type Entry interface {
	Cacheable

	Key() string
	Value() Value
}

type EvictionFn func(Entry)

type Engine interface {
	Insert(k string, v Value) bool
	Access(k string) (Value, bool)
	Size() uint64
	Dump() []Entry
	OnEvict(fn EvictionFn)
}
