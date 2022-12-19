package cache

import (
	"fmt"
	"log"
	"sync"

	"geek/cache/proto"
	"geek/cache/singleflight"
)

// 接收 key --> 检查是否被缓存 --是--> 返回缓存值 (1)
//                |
//                |--否--> 是否应当从远程节点获取 --是--> 与远程节点交互 --> 返回缓存值 (2)
//                            |
//                            |--否--> 调用`回调函数`，获取值并添加到缓存 --> 返回缓存值 (3)

// Getter loads data from given key
type Getter interface {
	Get(key string) ([]byte, error)
}

// GetterFunc implements Getter interface
type GetterFunc func(key string) ([]byte, error)

// Get implements the Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// Group cache with namespace
type Group struct {
	name   string
	getter Getter
	mcache cache
	peers  PeerPicker
	loader *singleflight.Group
}

// NewGroup create a new Group instance
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		mcache: cache{cacheBytes: cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// GetGroup returns cache group with given name
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

// Get return the key value from cache group
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key required")
	}

	if v, ok := g.mcache.get(key); ok {
		log.Printf("cache hit key: %q", key)
		return v, nil
	}

	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	fn := func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.Pick(key); ok {
				value, err := g.getFromPeer(peer, key)
				if err == nil {
					return value, nil
				}
				log.Printf("load from peer failed: %v", err)
			}
		}

		return g.getLocal(key)
	}

	v, err := g.loader.Do(key, fn)
	if err != nil {
		return ByteView{}, err
	}
	return v.(ByteView), nil
}

// getFromPeer with HTTP
// func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
// 	bytes, err := peer.Get(g.name, key)
// 	if err != nil {
// 		return ByteView{}, err
// 	}
// 	return ByteView{b: bytes}, nil
// }

func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &proto.Request{
		Group: g.name,
		Key:   key,
	}
	res := &proto.Response{}
	if err := peer.Get(req, res); err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

func (g *Group) getLocal(key string) (ByteView, error) {
	v, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	value := ByteView{cloneBytes(v)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mcache.add(key, value)
}

// RegisterPeers register a PeerPicker
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		// panic("RegisterPeers invoked more than once")
		return
	}
	g.peers = peers
}
