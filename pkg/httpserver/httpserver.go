package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"log"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

type Option func(*Server)

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

const (
	defaultAddr            = ":80"
	defaultMaxHeaderBytes  = 1 << 20
	defaultShutdownTimeout = 3 * time.Second
)

func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Addr:           defaultAddr,
		Handler:        handler,
		MaxHeaderBytes: defaultMaxHeaderBytes,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}
	s.start()

	return s
}

func (s *Server) start() {
	log.Printf("Starting HTTP server on port %s", s.server.Addr)
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
