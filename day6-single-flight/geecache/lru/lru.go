package lru

import "container/list"

type Value interface {
	Len() int
}
type entry struct {
	key   string
	value Value
}

type Cache struct {
	maxBytes  int64
	nBytes    int64
	ll        *list.List
	cache     map[string]*list.Element
	OnEvicted func(key string, value Value)
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
func NewCache(maxBytes int64, OnEvicted func(key string, value Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: OnEvicted,
	}
}

// Get 获取缓存
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

//Add 增加缓存
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		// 存在覆盖原值
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		kv.value = value
		c.nBytes += int64(value.Len() - kv.value.Len())
	} else {
		//不存在建立元素
		ele := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = ele
		kv := ele.Value.(*entry)
		c.nBytes += int64(len(key) + kv.value.Len())
	}
	for c.nBytes > c.maxBytes {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key) + kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}
