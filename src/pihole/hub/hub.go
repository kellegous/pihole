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
	glog.Infof("Unregistered: %s as %s", hd.ID(), host)

	delete(h.routes, host)
}

func (h *Hub) get(host string) http.Handler {
	h.lck.RLock()
	defer h.lck.RUnlock()
	return h.routes[host]
}

type response struct {
	http.ResponseWriter
	status int
}

func (r *response) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *response) log(req *http.Request) {
	glog.Infof("addr=%s code=%d method=%s host=%s uri=%s",
		req.RemoteAddr,
		r.status,
		req.Method,
		req.Host,
		req.RequestURI)
}

// ServeHTTP ...
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	res := response{
		ResponseWriter: w,
	}
	defer res.log(r)

	rt := h.get(r.Host)
	if rt == nil {
		http.NotFound(&res, r)
		return
	}
	rt.ServeHTTP(&res, r)
}
