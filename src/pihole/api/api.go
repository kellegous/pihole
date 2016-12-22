package api

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"pihole/hub"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// DefaultAddr ...
const DefaultAddr = "localhost:5000"

type api struct {
	h *hub.Hub
}

type proxy struct {
	p  *httputil.ReverseProxy
	c  chan string
	id string
}

func (c *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.p.ServeHTTP(w, r)
}

func (c *proxy) ID() string {
	return c.id
}

func (a *api) Ping(ctx context.Context, req *PingReq) (*PingRes, error) {
	return &PingRes{
		Id: req.Id,
	}, nil
}

func (a *api) Register(s Api_RegisterServer) error {
	r, err := s.Recv()
	if err != nil {
		return err
	}

	u, err := url.Parse(fmt.Sprintf("http://%s/", r.Addr))
	if err != nil {
		return err
	}

	c := proxy{
		p:  httputil.NewSingleHostReverseProxy(u),
		c:  make(chan string, 10),
		id: r.Id,
	}

	for _, host := range r.Hosts {
		a.h.Register(host, &c)
		defer a.h.Unregister(host, &c)
	}

	ch := make(chan error, 1)

	go func() {
		_, er := s.Recv()
		ch <- er
	}()

	for {
		select {
		case msg := <-c.c:
			if err := s.Send(&RegisterRes{
				Message: msg,
			}); err != nil {
				return err
			}
		case err = <-ch:
			return err
		}
	}
}

// ListenAndServe ...
func ListenAndServe(addr string, h *hub.Hub) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	defer l.Close()

	s := grpc.NewServer()
	RegisterApiServer(s, &api{h: h})

	glog.Infof("Api: %s", addr)
	return s.Serve(l)
}

// Dial ...
func Dial(c net.Conn) (*grpc.ClientConn, error) {
	return grpc.Dial(
		"tunnel",
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, t time.Duration) (net.Conn, error) {
			return c, nil
		}))
}
