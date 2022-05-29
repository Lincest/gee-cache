package geecache

import (
	"fmt"
	"reflect"
	"testing"
)

/**
    gee_cache
    @author: roccoshi
    @desc: test of lru
**/

type String string

// implements Len() of Value to get bytes
func (d String) Len() int {
	return len(d)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("k1", String("11"))
	if v, ok := lru.Get("k1"); !ok || string(v.(String)) != "11" {
		t.Fatalf("fail ([k1 => 11] map failed)")
	}
	if _, ok := lru.Get("k2"); ok {
		t.Fatalf("fail (get missing  key)")
	}
}

func TestExceed(t *testing.T) {
	removedKeys := make([]string, 0)
	cb := func(key string, value Value) {
		fmt.Printf("remove key = %s\n", key)
		removedKeys = append(removedKeys, key)
	}
	lru := New(int64(10), cb)
	lru.Add("k1", String("12345678")) // add size = 10 [k1]
	lru.Add("k2", String("1"))        // add size = 3 [k1]
	lru.Add("k3", String("1"))        // add size = 6 [k1]
	lru.Add("k4", String("1"))        // add size = 9 [k1]
	lru.Add("k5", String("1"))        // add size = 9 [k1, k2]
	except := []string{"k1", "k2"}
	if !reflect.DeepEqual(except, removedKeys) {
		t.Fatalf("except removedKeys equals to %s, current = %s", except, removedKeys)
	}
}
