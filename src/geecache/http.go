package geecache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

/**
    gee_cache
    @author: roccoshi
    @desc: http server
**/

const defaultBasePath = "/_geecache/"

// HTTPPool implements PeerPicker for a pool of HTTP peers (wtf)
type HTTPPool struct {
	self     string // e.g. https://localhost:1080
	basePath string // suffix, self + basePath is the real path
}

// NewHTTPPool is the constructor of HTTPPool
func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

// Log info with server name
// [Server Method] Messages
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Server %s] %s", p.self, fmt.Sprintf(format, v...))
}

// ServeHTTP is the main Server
func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// valid the path
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		http.Error(w, "unexpect server path", http.StatusBadRequest)
	}
	p.Log("%s, %s", r.Method, r.URL.Path)
	// basepath/groupname/key
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group: "+groupName, http.StatusNotFound)
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(view.ByteSlice()) // response body with value of key
}
