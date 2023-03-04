package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	//hash 函数
	hash Hash
	// 虚拟节点数量
	replicas int
	// 哈希环
	keys []int
	// 虚拟节点与真实节点的映射表
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	m := &Map{
		hash:     fn,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加节点
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			//获取虚拟节点hash
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			//加入哈希环
			m.keys = append(m.keys, hash)
			//加入虚拟节点与真实节点的映射
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 获取节点
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	//idx 的区间为0 -len(m.keys)
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]

}
