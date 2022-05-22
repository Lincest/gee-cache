package gee_cache

import "container/list"

/**
    gee_cache
    @author: roccoshi
    @desc: the go implementation of LRU(least recently used)

	two main data structures:
	1. 	map (hash map)
	2. 	double linked list in order to easier to move item and
		add / delete in O(1)
**/

// Cache (not safe for concurrent)
type Cache struct {
	maxBytes int64 // max memory
	nBytes int64 // used memory
	ll *list.List // double Linked List, Front is the newest and Back is the oldest
	cache map[string]*list.Element
	// optional life circle function: execute when entry is purged
	OnEvicted func(key string, value Value)
}

// Linked List's element
type entry struct {
	key string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New Cache's Constructor
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache {
		maxBytes: maxBytes,
		ll: list.New(),
		cache: make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get key -> value
func (c *Cache) Get(key string) (value Value, ok bool) {
	if e, ok := c.cache[key]; ok {
		c.ll.MoveToFront(e) // last used
		kv := e.Value.(*entry) // Value is an interface and coercion to entry
		return kv.value, true
	}
	// default ok is false
	return
}

// RemoveOldest removes the oldest item (which in the Back of list)
func (c *Cache) RemoveOldest() {
	e := c.ll.Back()
	if e != nil {
		c.ll.Remove(e)
		kv := e.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len()) // cause entry is key + value so memory is len(key) + size(value)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// Add to add a value to cache
func (c *Cache) Add(key string, value Value) {
	if e, ok := c.cache[key]; ok {
		// when element in cache
		c.ll.MoveToFront(e)
		kv := e.Value.(*entry)
		c.nBytes += int64(value.Len()) - int64(kv.value.Len()) // new - old
		kv.value = value // set to new value
	} else {
		// add to the front
		e := c.ll.PushFront(&entry{key, value})
		c.cache[key] = e
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// if memory exceed then remove oldest (use while because we don't know the size of oldest element)
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

// Len number of elements
func (c *Cache) Len() int {
	return c.ll.Len()
}