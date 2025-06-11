package server

import (
	"net"
	"net/http"
	"time"
)

type Server struct {
	listener net.Listener
	srv      *http.Server
}

func NewServer(handler http.Handler) (*Server, error) {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		return nil, err
	}
	srv := &http.Server{
		Handler:           handler,
		ReadHeaderTimeout: 30 * time.Second,
		ReadTimeout:       30 * time.Second,
	}
	return &Server{
		srv:      srv,
		listener: listener,
	}, nil
}

func (s *Server) Serve() error {
	if err := s.srv.Serve(s.listener); err != http.ErrServerClosed {
		return err
	}
	return nil
}
