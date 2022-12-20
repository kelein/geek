package cache

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"

	"geek/cache/consistenthash"
	"geek/cache/proto"

	probuf "github.com/golang/protobuf/proto"
)

const (
	defaultBasePath = "/_gcache/"
	defaultReplicas = 50
)

// * 检查 HTTPPool 是否实现 PeerPicker
var _ PeerPicker = (*HTTPPool)(nil)

// HTTPPool pool for HTTP peers
type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*httpGetter
}

// NewHTTPPool create HTTP pool of peers
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Set update HTTP pool of peers
func (p *HTTPPool) Set(peers ...string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.peers = consistenthash.New(defaultReplicas, nil)
	p.peers.Add(peers...)
	p.httpGetters = make(map[string]*httpGetter, len(peers))
	for _, peer := range peers {
		p.httpGetters[peer] = &httpGetter{
			baseURL: peer + p.basePath,
		}
	}
}

// Pick get a peer by the given key
func (p *HTTPPool) Pick(key string) (PeerGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	peer := p.peers.Get(key)
	if peer != "" && peer != p.self {
		p.Log("picked peer %q", peer)
		return p.httpGetters[peer], true
	}
	return nil, false
}

// Log record server info with name
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if !strings.HasPrefix(r.URL.Path, p.basePath) {
// 		panic("HTTPPool meets unexpected path: " + r.URL.Path)
// 	}
// 	p.Log("%s %s", r.Method, r.URL.Path)

// 	// * /<basePath>/<groupName>/<key> required
// 	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
// 	if len(parts) != 2 {
// 		http.Error(w, "bad request", http.StatusBadRequest)
// 		return
// 	}

// 	groupName := parts[0]
// 	group := GetGroup(groupName)
// 	if group == nil {
// 		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
// 		return
// 	}

// 	key := parts[1]
// 	view, err := group.Get(key)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	w.Header().Set("Content-Type", "application/octed-stream")
// 	w.Write(view.ByteSlice())
// }

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// * 检查 httpGetter 是否实现 PeerGetter
var _ PeerGetter = (*httpGetter)(nil)

type httpGetter struct {
	baseURL string
}

func (h *httpGetter) Get2(group, key string) ([]byte, error) {
	// format: http://example.com/_gcache/group/key
	link := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(group),
		url.QueryEscape(key),
	)
	res, err := http.Get(link)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server status %d", res.StatusCode)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("read response error: %v", err)
	}

	return bytes, nil
}

func (h *httpGetter) Get(in *proto.Request, out *proto.Response) error {
	// format: http://example.com/_gcache/group/key
	link := fmt.Sprintf(
		"%v%v/%v",
		h.baseURL,
		url.QueryEscape(in.GetGroup()),
		url.QueryEscape(in.GetKey()),
	)

	res, err := http.Get(link)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("server status %d", res.StatusCode)
	}

	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response body error: %v", err)
	}

	if err := probuf.Unmarshal(bytes, out); err != nil {
		return fmt.Errorf("unmarshal response error: %v", err)
	}

	return nil
}
