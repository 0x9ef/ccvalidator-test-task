package app

import (
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/0x9ef/card-validator/config"
	"github.com/0x9ef/card-validator/pkg/httpserver"
	"github.com/0x9ef/card-validator/pkg/logging"

	httpcontroller "github.com/0x9ef/card-validator/internal/controller/http"
	"github.com/0x9ef/card-validator/internal/service"
)

func Run(cfg *config.Config) {
	logger := logging.NewZap(cfg.Log.Level)
	rand.Seed(time.Now().UnixNano())

	services := service.Services{
		CardValidatorService: service.NewCreditCardValidator(cfg),
	}

	// Init router and handler
	httpHandler := gin.New()
	httpcontroller.New(&httpcontroller.Options{
		Handler:  httpHandler,
		Services: services,
		Logger:   logger,
		Config:   cfg,
	})

	// Init and run HTTP server
	httpServer := httpserver.New(httpHandler, httpserver.Port(cfg.HTTP.Port))

	// Waiting for a signal to do graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		logger.Info("app - Run - signal: " + s.String())
	case err := <-httpServer.Notify():
		logger.Error("app - Run - httpServer.Notify", "err", err)
	}

	if err := httpServer.Shutdown(); err != nil {
		logger.Error("app - Run - httpServer.Shutdown", "err", err)
	}
}
