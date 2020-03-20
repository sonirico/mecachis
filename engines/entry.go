package engines

type entry struct {
	key   string
	value Value
}

func NewEntry(k string, v Value) *entry {
	return &entry{key: k, value: v}
}

func (e *entry) Key() string {
	return e.key
}

func (e *entry) Value() Value {
	return e.value
}

func (e *entry) Len() uint64 {
	return uint64(len(e.key)) + e.value.Len()
}
