package consistenthash

import (
	"strconv"
	"testing"
)

func TestHash(t *testing.T) {
	keys := []string{"1", "2", "3"}
	fn := func(key []byte) uint32 {
		v, _ := strconv.Atoi(string(key))
		return uint32(v)
	}

	hashMap := New(3, fn)
	hashMap.Add(keys...)
	t.Logf("hashMap: %+v", hashMap)
	for _, k := range keys {
		v := hashMap.Get(k)
		t.Logf("hashMap.Get() key %q = %q", k, v)
	}
}
