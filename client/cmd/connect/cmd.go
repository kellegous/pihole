package connect

import (
	"time"

	"github.com/kellegous/pihole/client/config"
	"github.com/kellegous/pihole/client/proxy"

	"go.uber.org/zap"
)

// Name ...
const Name = "connect"

const reconnectDelay = 10 * time.Second

// Main ...
func Main(conf string, args []string) {
	var cfg config.Config
	if err := cfg.Read(conf); err != nil {
		zap.L().Fatal("unable to read config",
			zap.String("filename", conf),
			zap.Error(err))
	}

	for {
		if err := proxy.ConnectAndServeConfig(&cfg); err != nil {
			zap.L().Error("client error",
				zap.Error(err))
		}

		zap.L().Info("reconnecting",
			zap.Duration("delay", reconnectDelay),
			zap.Time("at", time.Now().Add(reconnectDelay)))
		time.Sleep(reconnectDelay)
	}
}
