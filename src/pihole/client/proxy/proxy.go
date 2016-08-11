package proxy

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"

	"pihole/api"
	"pihole/client/config"

	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

func keysFrom(cfg *config.Config) (ssh.AuthMethod, error) {
	prv, err := ioutil.ReadFile(cfg.PrivateKey())
	if err != nil {
		return nil, err
	}

	sgn, err := ssh.ParsePrivateKey(prv)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(sgn), nil
}

func register(c *ssh.Client, cfg *config.Config, addr string) error {
	t, err := c.Dial("tcp", api.DefaultAddr)
	if err != nil {
		return err
	}
	defer t.Close()

	cc, err := api.Dial(t)
	if err != nil {
		return err
	}
	defer cc.Close()

	ac := api.NewApiClient(cc)

	s, err := ac.Register(context.Background())
	if err != nil {
		return err
	}

	if err := s.Send(&api.RegisterReq{
		Hosts: cfg.Hosts,
		Addr:  addr,
	}); err != nil {
		return err
	}

	for {
		m, err := s.Recv()
		if err != nil {
			return err
		}

		log.Println(m)
	}
}

// ConnectAndServe ...
func ConnectAndServe(cfg *config.Config) error {
	auth, err := keysFrom(cfg)
	if err != nil {
		return err
	}

	c, err := ssh.Dial("tcp", cfg.Hub.Addr, &ssh.ClientConfig{
		User: cfg.Hub.User,
		Auth: []ssh.AuthMethod{
			auth,
		}})
	if err != nil {
		return err
	}
	defer c.Close()

	log.Printf("connected: %s", cfg.Hub.Addr)

	l, err := c.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}

	ch := make(chan error, 1)

	go func() {
		srv := http.Server{
			Handler: httputil.NewSingleHostReverseProxy(cfg.ToURL),
		}

		ch <- srv.Serve(l)
	}()

	go func() {
		ch <- register(c, cfg, l.Addr().String())
	}()

	return <-ch
}
