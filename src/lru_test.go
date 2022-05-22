package gee_cache

import "testing"

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