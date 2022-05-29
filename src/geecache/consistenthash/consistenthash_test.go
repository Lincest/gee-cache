package consistenthash

import (
	"strconv"
	"testing"
)

/**
    consistenthash
    @author: roccoshi
    @desc: test
**/

func TestHashing(t *testing.T) {
	// replica: x, 1x, 2x
	hash := New(3, func(key []byte) uint32 {
		i, _ := strconv.Atoi(string(key))
		return uint32(i)
	})
	hash.Add("6", "4", "2")
	testCases := map[string]string{
		"2":      "2",
		"11":     "2",
		"23":     "4",
		"27":     "2",
		"29":     "2",
		"100035": "2",
	}
	for k, v := range testCases {
		realV := hash.Get(k)
		if realV != v {
			t.Errorf("error, should be %s:%s but now %s:%s", k, v, k, realV)
		}
	}
	hash.Add("9")
	testCases["27"] = "9"
	testCases["29"] = "9"
	for k, v := range testCases {
		realV := hash.Get(k)
		if realV != v {
			t.Errorf("error, should be %s:%s but now %s:%s", k, v, k, realV)
		}
	}
}
