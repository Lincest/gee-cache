package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

/**
    consistenthash
    @author: roccoshi
    @desc: 一致性哈希
**/

// Hash bytes => uint32
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int            // 虚拟节点倍数
	keys     []int          // sorted hash circle
	hashMap  map[int]string // virtual node => real node
}

// New is the constructor
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // set default hash function, and approve user to select func
	}
	return m
}

// Add adds some keys to hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			// NODE: hash(i + real) = virtual
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// append virtual node
			m.keys = append(m.keys, hash)
			// virtual => real
			m.hashMap[hash] = key
		}
		sort.Ints(m.keys)
	}
}

// Get gets the closest item in the hash to the provided key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// binary search (if hash >= max(m.keys), then i = len(m.keys))
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
