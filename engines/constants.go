package engines

type CacheType int

const (
	LRU = iota
	LFU
	LFRU
	MRU
)

var cacheTypes = map[string]CacheType{
	"lru": LRU,
}

func LookupCacheType(candidate string) (CacheType, bool) {
	ct, ok := cacheTypes[candidate]
	return ct, ok
}
