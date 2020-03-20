package engines

type CacheType int

const (
	LRU = iota
	LFU
	LFRU
	MRU
)
