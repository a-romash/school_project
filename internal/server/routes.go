package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) SetRoutes(serveMux *mux.Router) {
	// serveMux.HandleFunc("/api", s.HandleAPIRoutes)
	serveMux.HandleFunc("/api/ping", s.HandlePing).Methods("GET")
	serveMux.HandleFunc("/api/login", s.HandleAPILogin).Methods("POST")
	serveMux.HandleFunc("/api/register", s.HandleAPIRegister).Methods("POST")

	serveMux.HandleFunc("/register", s.HandleRegister)
	serveMux.HandleFunc("/login", s.HandleLogin)
	serveMux.HandleFunc("/", s.HandleMain)

	serveMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
}
