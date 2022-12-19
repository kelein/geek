package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

// Map contains all hashed keys
type Map struct {
	keys     []int          // 哈希环
	replicas int            // 虚拟节点数量
	hashFunc Hash           // 哈希函数
	hashMap  map[int]string // 虚拟节点与真实节点映射
}

// New create a Map instance
func New(replicas int, fn Hash) *Map {
	m := &Map{
		hashFunc: fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hashFunc == nil {
		m.hashFunc = crc32.ChecksumIEEE
	}
	return m
}

// Add adds keys into to the map
func (m *Map) Add(keys ...string) {
	for _, k := range keys {
		for i := 0; i < m.replicas; i++ {
			// * 虚拟节点：strconv.Itoa(i) + key
			hash := int(m.hashFunc([]byte(strconv.Itoa(i) + k)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = k
		}
	}
	sort.Ints(m.keys)
}

// Get returns the nearest item by key
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	hash := int(m.hashFunc([]byte(key)))
	idx := sort.Search(len(m.keys), func(i int) bool {
		// Binary Search
		return m.keys[i] >= hash
	})

	// * idx 的值在 [0, n] 之间, 如果 idx=n
	// * 应该从哈希环头部取 m.keys[0]
	// * 故这里用查找key用idx取余的方式
	return m.hashMap[m.keys[idx%len(m.keys)]]
}

// Remove delete key from hash map
func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		hash := int(m.hashFunc([]byte(strconv.Itoa(i) + key)))
		idx := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:idx], m.keys[idx+1:]...)
		delete(m.hashMap, hash)
	}
}
