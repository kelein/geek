package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"geek/cache"
)

var (
	port = flag.Int("port", 8001, "cache server port")
	api  = flag.Bool("api", false, "enable api server")
)

var db = map[string]string{
	"Tom":  "630",
	"Jack": "550",
	"Sam":  "567",
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func example() {
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

func createGroup() *cache.Group {
	fn := func(key string) ([]byte, error) {
		log.Printf("[SlowDB] search key %q", key)
		if v, ok := db[key]; ok {
			return []byte(v), nil
		}
		return nil, fmt.Errorf("key %q not found", key)
	}

	return cache.NewGroup("scores", 2<<10, cache.GetterFunc(fn))
}

func startCacheServer(addr string, addrs []string, group *cache.Group) {
	peers := cache.NewHTTPPool(addr)
	peers.Set(addrs...)
	group.RegisterPeers(peers)
	log.Printf("cache is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr[7:], peers))
}

func startAPIServer(addr string, group *cache.Group) {
	api := func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		view, err := group.Get(key)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/octed-stream")
		w.Write(view.ByteSlice())
	}

	http.Handle("/api", http.HandlerFunc(api))
	log.Printf("api server is running at %s", addr)
	log.Fatal(http.ListenAndServe(addr[7:], nil))
}

func main() {
	flag.Parse()

	apiAddr := "http://localhost:9000"
	addrMap := map[int]string{
		8001: "http://localhost:8001",
		8002: "http://localhost:8002",
		8003: "http://localhost:8003",
	}
	addrs := []string{}
	for _, v := range addrMap {
		addrs = append(addrs, v)
	}

	cacher := createGroup()
	if *api {
		go startAPIServer(apiAddr, cacher)
	}
	startCacheServer(addrMap[*port], addrs, cacher)
}
