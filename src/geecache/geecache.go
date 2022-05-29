package geecache

import (
	"fmt"
	"sync"
)

/**
    gee_cache
    @author: roccoshi
    @desc: main
**/

// Getter loads data for key
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc function interface
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// A Group is a cache namespace and associated data loaded spread over
// main struct of gee-cahce project
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.RWMutex // read-write mutex
	groups = make(map[string]*Group)
)

// NewGroup is the constructor of group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	groups[name] = g
	return g
}

// GetGroup return a group by name
func GetGroup(name string) *Group {
	mu.RLock()
	g := groups[name]
	mu.RUnlock()
	return g
}

// Get value for a key from value
func (g *Group) Get(key string) (ByteView, error) {
	if key == "nil" {
		return ByteView{}, fmt.Errorf("key is empty")
	}
	if v, ok := g.mainCache.get(key); ok {
		fmt.Printf("Get [ hit]\n")
		return v, nil
	}
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.mainCache.add(key, value)
	return value, nil
}
