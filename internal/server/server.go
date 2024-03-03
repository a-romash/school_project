package server

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	address string
	logger  *slog.Logger
}

func NewServer(address string, logger *slog.Logger) *Server {
	return &Server{
		address: address,
		logger:  logger,
	}
}

func (s *Server) Start() error {
	serveMux := mux.NewRouter()

	s.SetRoutes(serveMux)

	s.logger.Info("server has been started", "address", s.address)

	err := http.ListenAndServe(s.address, serveMux)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
