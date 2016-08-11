package hub

import (
	"log"
	"net/http"
	"sync"
)

// Hub ...
type Hub struct {
	routes map[string]http.Handler
	lck    sync.RWMutex
}

// NewHub ...
func NewHub() *Hub {
	return &Hub{
		routes: map[string]http.Handler{},
	}
}

// Register ...
func (h *Hub) Register(host string, hd http.Handler) error {
	h.lck.Lock()
	defer h.lck.Unlock()

	h.routes[host] = hd

	return nil
}

// Unregister ...
func (h *Hub) Unregister(host string) {
	h.lck.Lock()
	defer h.lck.Unlock()
	log.Printf("unregister %s", host)
	delete(h.routes, host)
}

func (h *Hub) get(host string) http.Handler {
	h.lck.RLock()
	defer h.lck.RUnlock()
	return h.routes[host]
}

// ServeHTTP ...
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.String())
	rt := h.get(r.Host)
	if rt == nil {
		http.NotFound(w, r)
		return
	}
	rt.ServeHTTP(w, r)
}
