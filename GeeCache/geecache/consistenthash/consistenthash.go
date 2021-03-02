package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

type Map struct {
	hash         Hash           // Hash function
	replicas     int            // num of replicated virtual nodes
	virtualNodes []int          // sorted virtualNodes
	hashMap      map[int]string // virtual node -> real node
}

// New creates a Map instance
func New(replicas int, hashfunc Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     hashfunc,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add some keys(i.e. real nodes) to the hash
func (m *Map) Add(realNodes ...string) {
	for _, node := range realNodes {
		for i := 0; i < m.replicas; i++ {
			virtualNode := int(m.hash([]byte(strconv.Itoa(i) + node)))
			m.virtualNodes = append(m.virtualNodes, virtualNode)
			m.hashMap[virtualNode] = node
		}
		sort.Ints(m.virtualNodes)
	}
}

func (m *Map) Get(key string) string {
	if len(m.virtualNodes) == 0 {
		return ""
	}

	hashValue := int(m.hash([]byte(key)))
	// Binary search for approriate replica
	idx := sort.Search(len(m.virtualNodes), func(i int) bool {
		return m.virtualNodes[i] >= hashValue
	})

	return m.hashMap[m.virtualNodes[idx%len(m.virtualNodes)]] // return the real key
}
