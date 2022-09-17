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
	chain  *list.List
	cache  map[string]*list.Element

	OnEvicted func(key string, value Value)
}

// New create a Cache instance with capacity
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		chain:     list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get find the key value with exist state
func (c *Cache) Get(key string) (Value, bool) {
	if item, ok := c.cache[key]; ok {
		c.chain.MoveToFront(item)
		e := item.Value.(*entry)
		return e.value, true
	}
	return nil, false
}

// RemoveOld removes the oldest item from cache
func (c *Cache) RemoveOld() {
	item := c.chain.Back()
	if item != nil {
		c.chain.Remove(item)
		e := item.Value.(*entry)
		delete(c.cache, e.key)
		c.nbytes -= int64(len(e.key) + e.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(e.key, e.value)
		}
	}
}

// Add adds the given value to the cache
func (c *Cache) Add(key string, value Value) {
	if item, ok := c.cache[key]; ok {
		c.chain.MoveToFront(item)
		e := item.Value.(*entry)
		e.value = value
		c.nbytes += int64(value.Len() - e.value.Len())
	} else {
		item := c.chain.PushFront(&entry{key, value})
		c.cache[key] = item
		c.nbytes += int64(len(key) + value.Len())
	}

	// Remove oldest items when beyond the limit
	for c.maxBytes != 0 && c.nbytes > c.maxBytes {
		c.RemoveOld()
	}
}

// Len returns items number in the cache
func (c *Cache) Len() int {
	return c.chain.Len()
}
