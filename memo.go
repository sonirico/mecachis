package mecachis

type MemoryView []byte

func (mv MemoryView) Value() interface{} {
	res := make([]byte, len(mv))
	copy(res, mv)
	return res
}

func (mv MemoryView) Len() uint64 {
	return uint64(len(mv))
}

func (mv MemoryView) String() string {
	return string(mv)
}

func (mv MemoryView) Clone() []byte {
	res := make([]byte, len(mv))
	copy(res, mv)
	return res
}
