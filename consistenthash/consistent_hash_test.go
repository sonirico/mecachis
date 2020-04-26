package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

type testS struct {
	key  string
	node string
}

func testHasher(data []byte) int {
	i, err := strconv.Atoi(string(data))
	if err != nil {
		panic(err)
	}
	return i
}

func testConsistentMap(t *testing.T, c *ConsistentHash, tests []testS) {
	t.Helper()

	for _, test := range tests {
		if c.Get(test.key) != test.node {
			t.Errorf("unexpected node for key %s. want %s, have %s",
				test.key, test.node, c.Get(test.key))
		}
	}
}

func TestConsistentHash_Add(t *testing.T) {
	keys := []string{"1", "2", "3"} // [11, 21, 12, 22, 13, 23] -> // [11, 12, 13, 21, 21, 22, 23]
	cmap := New(2, testHasher)
	cmap.Add(keys...)
	tests := []testS{
		{"3", "1"},
		{"12", "2"},
		{"13", "3"},
		{"10", "1"},
		{"21", "1"},
		{"22", "2"},
		{"23", "3"},
		{"24", "1"},
	}

	testConsistentMap(t, cmap, tests)

	cmap.Add("8") // [11, 21, 12, 22, 13, 23, 18, 81] -> // [11, 12, 13, 18, 21, 21, 22, 23, 81]
	tests = tests[:len(tests)-1]
	tests = append(tests, testS{"24", "8"})
	tests = append(tests, testS{"25", "8"})

	testConsistentMap(t, cmap, tests)

}

func BenchmarkGet8_MD5(b *testing.B)     { benchmarkGet(b, 8, HashMD5) }
func BenchmarkGet32_MD5(b *testing.B)    { benchmarkGet(b, 32, HashMD5) }
func BenchmarkGet128_MD5(b *testing.B)   { benchmarkGet(b, 128, HashMD5) }
func BenchmarkGet512_MD5(b *testing.B)   { benchmarkGet(b, 512, HashMD5) }
func BenchmarkGet8_CRC32(b *testing.B)   { benchmarkGet(b, 8, HashCRC32) }
func BenchmarkGet32_CRC32(b *testing.B)  { benchmarkGet(b, 32, HashCRC32) }
func BenchmarkGet128_CRC32(b *testing.B) { benchmarkGet(b, 128, HashCRC32) }
func BenchmarkGet512_CRC32(b *testing.B) { benchmarkGet(b, 512, HashCRC32) }

func benchmarkGet(b *testing.B, shards int, fn HasherFunc) {

	hash := New(50, fn)

	var buckets []string
	for i := 0; i < shards; i++ {
		buckets = append(buckets, fmt.Sprintf("shard-%d", i))
	}

	hash.Add(buckets...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash.Get(buckets[i&(shards-1)])
	}
}
