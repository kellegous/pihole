package hub

import (
	"net/http"
	"sync"

	"go.uber.org/zap"
)

// Proxy ...
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

	zap.L().Info("registered",
		zap.String("id", hd.ID()),
		zap.String("host", host))
	h.routes[host] = hd

	return nil
}

// Unregister ...
func (h *Hub) Unregister(host string, exp Proxy) {
	h.lck.Lock()
	defer h.lck.Unlock()

	hd := h.routes[host]
	if hd != exp {
		zap.L().Warn("unregister with non-matching proxy",
			zap.String("host", host))
		return
	}

	zap.L().Info("unregistered",
		zap.String("id", hd.ID()),
		zap.String("host", host))

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
