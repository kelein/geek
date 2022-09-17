package lru

import (
	"reflect"
	"testing"
)

type StrEntry string

func (s StrEntry) Len() int { return len(s) }

func TestCache_Get(t *testing.T) {
	lru := New(0, nil)

	type args struct{ key string }

	tests := []struct {
		name  string
		args  args
		want  Value
		want1 bool
	}{
		{"A", args{"AAA"}, StrEntry("aaa"), true},
		{"B", args{"BBB"}, StrEntry("bbb"), true},
		{"C", args{"CCC"}, StrEntry("ccc"), true},
		{"D", args{"DDD"}, StrEntry("ddd"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lru.Add(tt.args.key, tt.want)
		})
	}
	t.Logf("current cache size: %v", lru.Len())

	// * Test add items already in the cache
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lru.Add(tt.args.key, tt.want)
		})
	}
	t.Logf("current cache size: %v", lru.Len())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1 := lru.Get(tt.args.key)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCache_RemoveOld(t *testing.T) {
	type args struct{ key string }

	tests := []struct {
		name  string
		args  args
		want  Value
		want1 bool
	}{
		{"A", args{"AAA"}, StrEntry("aaa"), false},
		{"B", args{"BBB"}, StrEntry("bbb"), false},
		{"C", args{"CCC"}, StrEntry("ccc"), true},
		{"D", args{"DDD"}, StrEntry("ddd"), true},
	}

	size := 0
	for _, e := range tests[:2] {
		size += len(e.args.key) + e.want.Len()
	}

	onEvicted := func(k string, v Value) {
		t.Logf("evicted entry <%s:%s>", k, v)
	}
	lru := New(int64(size), onEvicted)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lru.Add(tt.args.key, tt.want)
		})
	}

	t.Logf("current cache size: %v", lru.Len())

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := lru.Get(tt.args.key)
			if i > 2 && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Cache.Get(%s) got = %v, want %v", tt.args.key, got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Cache.Get(%s) got1 = %v, want %v", tt.args.key, got1, tt.want1)
			}
		})
	}
}
