package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) SetRoutes(serveMux *mux.Router) {
	serveMux.HandleFunc("/api/ping", s.HandlePing)
	serveMux.HandleFunc("/register", s.HandleRegister)
	serveMux.HandleFunc("/login", s.HandleLogin)
	serveMux.HandleFunc("/", s.HandleMain)

	serveMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
}
