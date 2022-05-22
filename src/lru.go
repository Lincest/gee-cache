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
	maxBytes int64
	nbytes int64
	ll *list.List
	cache map[string]*list.Element
	OnEvicted func(key string, value Value)
}

// TODO: what is this
type entry struct {
	key string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}