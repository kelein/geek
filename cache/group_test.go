package cache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetterFunc_Get(t *testing.T) {
	var g Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"A", args{"TestKey"}, []byte("TestKey"), false},
		{"B", args{"getterFunc"}, []byte("getterFunc"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := g.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetterFunc.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetterFunc.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

var mockDB = map[string]string{
	"TOM": "630",
	"JAN": "589",
	"SAM": "547",
	"MOY": "568",
}

var loadNum = make(map[string]int, len(mockDB))

var getScore = func(key string) ([]byte, error) {
	log.Printf("search key: %q", key)
	if v, ok := mockDB[key]; ok {
		if _, ok := loadNum[key]; !ok {
			loadNum[key] = 0
		}
		loadNum[key]++
		return []byte(v), nil
	}

	return nil, fmt.Errorf("%s not found", key)
}

var gcache = NewGroup("score", 2<<10, GetterFunc(getScore))

func TestGroup_Get(t *testing.T) {
	type testcase struct {
		name    string
		key     string
		want    ByteView
		wantErr bool
	}

	tests := []testcase{}
	for k, v := range mockDB {
		tcase := testcase{k, k, ByteView{[]byte(v)}, false}
		tests = append(tests, tcase)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := gcache.Get(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Group.Get(%q) error = %v, wantErr %v", tt.key, err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Group.Get(%q) = %v, want %v", tt.key, got, tt.want)
			}

			view, err := gcache.Get(tt.key)
			if err != nil {
				t.Errorf("Group.Get(%q) error = %v", tt.key, err)
			}
			if loadNum[tt.key] > 1 {
				t.Errorf("cache key %q miss", tt.key)
			}
			log.Printf("cache key %q value: %q", tt.key, view.String())
		})
	}

	log.Printf("cache key load count: %v", loadNum)
}

func TestGetGroup(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want *Group
	}{
		{"A", args{"score"}, gcache},
		{"B", args{"unknow"}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetGroup(tt.args.name); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGroup() = %v, want %v", got, tt.want)
			}

			if _, err := gcache.Get(""); err == nil {
				t.Errorf("group.Get() key must validated")
			}

			v, err := gcache.Get("Someone")
			if err != nil {
				t.Logf("get not exist key %q err: %v", "Someone", err)
			}
			t.Logf("got not exist key %q value: %q", "Someone", v.String())
		})
	}
}
