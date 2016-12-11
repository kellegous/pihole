package hub

import (
	"net/http"
	"sync"

	"github.com/golang/glog"
)

type Proxy interface {
	http.Handler
	ID() string
}

// Hub ...
type Hub struct {
	routes map[string]Proxy
	lck    sync.RWMutex
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		routes: map[string]Proxy{},
	}
}

// Register ...
func (h *Hub) Register(host string, hd Proxy) error {
	h.lck.Lock()
	defer h.lck.Unlock()

	glog.Infof("Registered: %s as %s", hd.ID(), host)
	h.routes[host] = hd

	return nil
}

// Unregister ...
func (h *Hub) Unregister(host string) {
	h.lck.Lock()
	defer h.lck.Unlock()

	hd := h.routes[host]
	if hd == nil {
		glog.Warningf("Attempt to unregister %s, but not registered.", host)
		return
	}

	glog.Infof("Unregistered: %s as %s", hd.ID(), host)

	delete(h.routes, host)
}

func (h *Hub) get(host string) http.Handler {
	h.lck.RLock()
	defer h.lck.RUnlock()
	return h.routes[host]
}

// ServeHTTP ...
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rt := h.get(r.Host)
	if rt == nil {
		http.NotFound(w, r)
		return
	}
	rt.ServeHTTP(w, r)
}
