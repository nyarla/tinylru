package tinylru

import (
	"sync"
)

type item[K comparable, V any] struct {
	key   K
	value V
	front *item[K, V]
	back  *item[K, V]
}

type Cache[K comparable, V any] struct {
	first *item[K, V]
	last  *item[K, V]
	data  map[K]*item[K, V]
	len   int
	size  int
	pool  sync.Pool
}

func (c *Cache[K, V]) Init(size int) *Cache[K, V] {
	c.first = nil
	c.last = nil
	c.data = make(map[K]*item[K, V])
	c.size = size

	c.pool = sync.Pool{
		New: func() any {
			return new(item[K, V])
		},
	}

	c.len = 0

	return c
}

func (c *Cache[K, V]) move(entry *item[K, V]) {
	// entry is head item
	if entry.front == nil {
		return
	}

	// entry is last item
	if entry.front != nil && entry.back == nil {
		entry.front.back = nil
		c.last = entry.front
	}

	if entry.front != nil && entry.back != nil {
		entry.front.back = entry.back
		entry.back.front = entry.front
	}

	entry.front = nil

	entry.back = c.first
	c.first.front = entry

	c.first = entry
	c.data[c.first.key] = c.first
}

func (c *Cache[K, V]) Remove(key K) bool {
	if entry, ok := c.data[key]; ok {
		if entry.front != nil {
			entry.front.back = entry.back
		}

		if entry.back != nil {
			entry.back.front = entry.front
		}

		entry.front = nil
		entry.back = nil

		delete(c.data, key)

		c.pool.Put(entry)
		c.len--

		return true
	}

	return false
}

func (c *Cache[K, V]) Get(key K) (value V, ok bool) {
	entry, exist := c.data[key]
	if exist {
		value = entry.value
		ok = exist

		c.move(entry)
	}

	return
}

func (c *Cache[K, V]) Put(key K, value V) V {
	if entry, ok := c.data[key]; ok {
		c.move(entry)
		return entry.value
	}

	entry := c.pool.Get().(*item[K, V])
	entry.key = key
	entry.value = value

	c.data[key] = entry

	c.len++

	// first put
	if c.first == nil {
		c.first = entry
		return entry.value
	}

	// second put
	if c.last == nil {
		c.last = c.first

		entry.front = nil
		entry.back = c.first

		c.first.front = entry
		c.first = entry

		return c.first.value
	}

	// other put
	entry.front = nil
	entry.back = c.first

	c.first.front = entry
	c.first = entry

	if c.len > c.size {
		last := c.last
		c.last = c.last.front
		c.last.back = nil

		c.Remove(last.key)
	}

	return entry.value
}
