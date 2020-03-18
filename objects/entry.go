package objects

import mc "github.com/sonirico/mecachis"

type entry struct {
	key   string
	value mc.Value
}

func NewEntry(k string, v mc.Value) *entry {
	return &entry{key: k, value: v}
}

func (e *entry) Key() string {
	return e.key
}

func (e *entry) Value() mc.Value {
	return e.value
}

func (e *entry) Len() uint64 {
	return uint64(len(e.key) + e.value.Len())
}
