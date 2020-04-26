package consistenthash

import (
	"crypto/md5"
	"encoding/binary"
	"hash/crc32"
	"sort"
	"strconv"
)

type HasherFunc func([]byte) int

type Hasher interface {
	Hash() HasherFunc
}

type ConsistentHash struct {
	hasher   HasherFunc
	replicas int
	hashmap  map[int]string
	keys     []int
}

func New(replicas int, fn HasherFunc) *ConsistentHash {
	return &ConsistentHash{
		hasher:   fn,
		replicas: replicas,
		hashmap:  make(map[int]string),
	}
}

func (ch *ConsistentHash) Hash(data []byte) int {
	return ch.hasher(data)
}

func (ch *ConsistentHash) Add(keys ...string) *ConsistentHash {
	for _, key := range keys {
		for i := 1; i <= ch.replicas; i++ {
			hash := ch.Hash([]byte(strconv.Itoa(i) + key))
			ch.keys = append(ch.keys, hash)
			ch.hashmap[hash] = key
		}
	}
	sort.Ints(ch.keys)
	return ch
}

func (ch *ConsistentHash) Get(key string) string {
	hash := ch.hasher([]byte(key))
	klen := len(ch.keys)
	ki := sort.Search(klen, func(i int) bool {
		return ch.keys[i] >= hash
	})
	if ki == klen {
		ki = 0
	}
	return ch.hashmap[ch.keys[ki]]
}

func HashCRC32(data []byte) int {
	return int(crc32.ChecksumIEEE(data))
}

func HashMD5(data []byte) int {
	hashed := md5.Sum(data)
	slice := hashed[:]
	return int(binary.LittleEndian.Uint32(slice))
}
