package lru

import "container/list"

// Value for abstract
type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

// Cache with LRU
type Cache struct {
	// maxBytes for capacity
	maxBytes int64
	// current used bytes
	nbytes int64
	l      *list.List
	cache  map[string]*list.Element

	OnEvicted func(key string, value Value)
}

// New create a Cache instance with capacity
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		l:         list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get find the key value with exist state
func (c *Cache) Get(key string) (Value, bool) {
	if item, ok := c.cache[key]; ok {
		c.l.MoveToFront(item)
		e := item.Value.(*entry)
		return e.value, true
	}
	return nil, false
}
