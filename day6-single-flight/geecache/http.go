package geecache

import (
	"fmt"
	"github.com/ATao13/gee-chache/day6-single-flight/geecache/consistenthash"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type HTTPGetter struct {
	baseURL string
}

func (h *HTTPGetter) Get(group string, key string) ([]byte, error) {
	r := fmt.Sprintf("%s%s/%s", h.baseURL, url.QueryEscape(group), url.QueryEscape(key))
	resp, err := http.Get(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server return:%d", resp.StatusCode)
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body error:%v", err)
	}
	return bytes, nil
}

var _ PeerGetter = (*HTTPGetter)(nil)

const (
	defaultBasePath = "/_geecache/"
	defaltReplicas  = 50
)

type HTTPPool struct {
	self        string
	basePath    string
	mu          sync.Mutex
	peers       *consistenthash.Map
	httpGetters map[string]*HTTPGetter
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// 初始化哈希算法
//设置节点
//设置节点何http url 的映射关系
func (h *HTTPPool) Set(peers ...string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.peers = consistenthash.New(defaltReplicas, nil)
	h.peers.Add(peers...)
	h.httpGetters = make(map[string]*HTTPGetter, len(peers))
	for _, peer := range peers {
		h.httpGetters[peer] = &HTTPGetter{baseURL: peer + h.basePath}
	}
}
func (h *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[server %s] %s", h.self, fmt.Sprintln(format, v))
}

func (h *HTTPPool) PickPeer(key string) (PeerGetter, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if peer := h.peers.Get(key); peer != "" && peer != h.self {
		h.Log("picc peer:%s", peer)
		return h.httpGetters[peer], true
	}
	return nil, false
}
func (h *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, h.basePath) {
		panic("HTTPPool")
	}
	parts := strings.SplitN(r.URL.Path[len(h.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	group := groups[groupName]
	key := parts[1]
	if group == nil {
		http.Error(w, "no such group"+groupName, http.StatusNotFound)
		return
	}
	value, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(value.ByteSlice())
	return

}
