package http

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/0x9ef/card-validator/config"
	"github.com/0x9ef/card-validator/internal/service"
	"github.com/0x9ef/card-validator/pkg/logging"

	"github.com/DataDog/gostackparse"
	"github.com/gin-gonic/gin"
)

type ErrorType string

const (
	ErrorTypeClient ErrorType = "client"
	ErrorTypeServer ErrorType = "Server"
)

type ErrorCode string

const (
	ErrorCodePanic ErrorCode = "01H"
)

type Options struct {
	Handler  *gin.Engine
	Services service.Services
	Logger   logging.Logger
	Config   *config.Config
}

type RouterOptions struct {
	Handler  *gin.RouterGroup
	Services service.Services
	Logger   logging.Logger
	Config   *config.Config
}

type RouterContext struct {
	services service.Services
	logger   logging.Logger
}

func New(options *Options) {
	routerOptions := &RouterOptions{
		Handler:  options.Handler.Group("/api/v1"),
		Services: options.Services,
		Logger:   options.Logger.Named("HTTPController"),
		Config:   options.Config,
	}

	routerOptions.Handler.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })
	{
		newCreditCardValidatorRoutes(routerOptions)
	}
}

func errorHandlerMiddleware(options *RouterOptions, handler func(c *gin.Context) (interface{}, *httpError)) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := options.Logger.Named("errorHandler")

		defer func() {
			if err := recover(); err != nil {
				stacktrace := debug.Stack()
				fmt.Fprintf(os.Stdout, string(stacktrace))

				goroutines, errors := gostackparse.Parse(bytes.NewReader(stacktrace))
				if len(errors) > 0 || len(goroutines) == 0 {
					logger.Error("got stacktrace", "stacktrace", string(stacktrace), "stacktraceErrors", errors, "err", err)
				} else {
					logger.Error("unhandled error", "err", err, "stacktrace", string(stacktrace))
				}

				c.AbortWithStatusJSON(http.StatusInternalServerError, httpError{
					Type:    ErrorTypeClient,
					Code:    ErrorCodePanic,
					Message: fmt.Sprintf("%+v", err),
				})
			}
		}()

		// Execute handler and check if it is a middleware
		body, err := handler(c)
		if body == nil && err == nil {
			return
		}

		if err != nil && err.Type == ErrorTypeServer {
			c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		} else if err != nil && err.Type == ErrorTypeClient {
			c.AbortWithStatusJSON(err.StatusCode, err)
		} else {
			logger.Info("request handled successfully")
			c.JSON(http.StatusOK, body)
		}
	}
}
