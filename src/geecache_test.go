package gee_cache

import (
	"fmt"
	"log"
	"testing"
)

/**
    gee_cache
    @author: roccoshi
    @desc: test of geecache
**/

var db = map[string]string{
	"asd": "123",
	"fgh": "456",
	"jkl": "789",
}

func TestGroup_Get(t *testing.T) {
	// count the number of callback funcs' calls
	loadCounts := make(map[string]int, len(db))
	// define a func
	gee := NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[db search key]", key)
			if v, ok := db[key]; ok {
				if _, ok := loadCounts[key]; !ok {
					loadCounts[key] = 0
				}
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))
	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		}
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		}
	}
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of [unknown] shoud be empty, but %s got", view)
	}
}
