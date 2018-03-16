package getkey

import (
	"fmt"

	"go.uber.org/zap"

	"pihole/client/config"
)

// Name ...
const Name = "get-key"

// Main ...
func Main(conf string, args []string) {
	var cfg config.Config
	if err := cfg.Read(conf); err != nil {
		zap.L().Fatal("unable to read config",
			zap.String("filename", conf),
			zap.Error(err))
	}

	pub, err := cfg.PublicKey()
	if err != nil {
		zap.L().Fatal("unable to get public key",
			zap.Error(err))
	}

	fmt.Println(pub)
}
