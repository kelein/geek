package main

import (
	"fmt"
	"log"
	"net/http"

	"geek/cache"
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "550",
	"Sam":  "567",
}

func main() {
	fn := func(key string) ([]byte, error) {
		log.Printf("[SlowDB] search key %q", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %q not found", key)
	}

	cache.NewGroup("scores", 2<<10, cache.GetterFunc(fn))

	addr := "localhost:9990"
	peers := cache.NewHTTPPool(addr)
	log.Printf("cache service running on %s", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
