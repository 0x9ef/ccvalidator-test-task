package main

import (
	"github.com/0x9ef/card-validator/internal/app"

	"github.com/0x9ef/card-validator/config"
	"github.com/0x9ef/card-validator/pkg/logging"
)

func main() {
	// Init logger and config
	logger := logging.NewZap("main")
	cfg := config.Get()
	logger.Info("read config", "config", cfg)

	app.Run(cfg)
}
