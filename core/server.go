package core

import (
	"context"
	"net/http"
)

type ServerOption func(s *Server)

func ListenAddress(addr string) ServerOption {
	return func(s *Server) {
		s.listenAddress = addr
	}
}

type Server struct {
	server http.Server
	listenAddress string
	AppContext AppContext
}

func NewServer(appContext AppContext, handler http.Handler, opts ...ServerOption) *Server {
	s := Server{
		listenAddress: ":3000",
		AppContext: appContext,
	}

	for _, opt := range opts {
		opt(&s)
	}

	// construct the server
	mux := http.NewServeMux()
	mux.Handle("/", handler)
	s.server = http.Server{
		Addr: s.listenAddress,
		Handler: mux,
	}

	return &s
}

func (s *Server) ListenAndServe(ctx context.Context) error {
	err := s.AppContext.Node.Bootstrap(ctx)
	if err != nil {
		return err
	}

	return s.server.ListenAndServe()
}

func (s *Server) Addr() string {
	return s.server.Addr
}
