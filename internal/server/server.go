package server

import (
	"log/slog"
	"net/http"
	"project/pkg/database/postgresql"

	"github.com/gorilla/mux"
)

type Server struct {
	address string
	db      *postgresql.Postgresql
}

func NewServer(address string, db *postgresql.Postgresql) *Server {
	return &Server{
		address: address,
		db:      db,
	}
}

func (s *Server) Start() error {
	serveMux := mux.NewRouter()

	s.SetRoutes(serveMux)

	slog.Info("server has been started", "address", s.address)

	err := http.ListenAndServe(s.address, serveMux)
	if err != http.ErrServerClosed {
		return err
	}

	return nil
}
