package proxy

import (
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/kellegous/pihole/api"
	"github.com/kellegous/pihole/client/config"
	"github.com/kellegous/pihole/logging"

	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"golang.org/x/net/context"
)

func sshAuthMethod(key []byte) (ssh.AuthMethod, error) {
	sgn, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, err
	}

	return ssh.PublicKeys(sgn), nil
}

func register(opts *Options, c *ssh.Client, addr string) error {
	t, err := c.Dial("tcp", opts.APIAddr)
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
		Hosts: opts.ClientHosts,
		Addr:  addr,
		Id:    opts.ClientID,
	}); err != nil {
		return err
	}

	go func() {
		defer c.Close()
		defer cc.Close()

		for i := 1; ; i++ {
			ctx, _ := context.WithTimeout(
				context.Background(),
				10*time.Second)
			if _, err := ac.Ping(ctx, &api.PingReq{
				Id: int64(i),
			}); err != nil {
				zap.L().Error("ping failed",
					zap.Error(err))
				return
			}

			time.Sleep(5 * time.Second)
		}
	}()

	for {
		m, err := s.Recv()
		if err != nil {
			return err
		}

		zap.L().Info("message",
			zap.String("msg", m.String()))
	}
}

// ConnectAndServe ...
func ConnectAndServe(
	opts *Options,
	h http.Handler) error {
	auth, err := sshAuthMethod(opts.SSHPrivateKey)
	if err != nil {
		return err
	}

	c, err := ssh.Dial("tcp", opts.SSHAddr, &ssh.ClientConfig{
		User:    opts.SSHUser,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			zap.L().Info("host key callback",
				zap.String("host", hostname),
				zap.String("addr", remote.String()),
				zap.String("key-type", key.Type()))
			return nil
		},
		Auth: []ssh.AuthMethod{
			auth,
		}})
	if err != nil {
		return err
	}
	defer c.Close()

	zap.L().Info("connected",
		zap.String("to", opts.SSHAddr))

	l, err := c.Listen("tcp", "localhost:0")
	if err != nil {
		return err
	}
	defer l.Close()

	ch := make(chan error, 1)

	go func() {
		srv := http.Server{
			Handler: h,
		}

		ch <- srv.Serve(l)
	}()

	go func() {
		ch <- register(opts, c, l.Addr().String())
	}()

	return <-ch
}

// ConnectAndServeConfig ...
func ConnectAndServeConfig(cfg *config.Config) error {
	key, err := ioutil.ReadFile(cfg.PrivateKey())
	if err != nil {
		return err
	}

	return ConnectAndServe(&Options{
		SSHAddr:       cfg.Hub.Addr,
		SSHUser:       cfg.Hub.User,
		SSHPrivateKey: key,
		APIAddr:       api.DefaultAddr,
		ClientHosts:   cfg.Hosts,
		ClientID:      cfg.ID,
	}, logging.WithLog(
		httputil.NewSingleHostReverseProxy(cfg.ToURL)))
}
