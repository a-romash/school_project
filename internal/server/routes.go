package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) SetRoutes(serveMux *mux.Router) {
	// serveMux.HandleFunc("/api", s.HandleAPIRoutes)
	serveMux.HandleFunc("/api/v1/ping", s.HandlePing).Methods("GET")
	serveMux.HandleFunc("/api/v1/login", s.HandleAPILogin).Methods("POST")
	serveMux.HandleFunc("/api/v1/register", s.HandleAPIRegister).Methods("POST")
	serveMux.HandleFunc("/api/v1/createtest", s.HandleApiCreateNewTest).Methods("POST")
	serveMux.HandleFunc("/api/v1/gettest", s.HandleApiGetTest)
	serveMux.HandleFunc("/api/v1/getinfo", s.HandleApiGetInfo).Methods("POST")
	serveMux.HandleFunc("/api/v1/getresult", s.HandleApiGetResult).Methods("POST")
	serveMux.HandleFunc("/api/v1/deletetoken", s.HandleApiDeleteToken).Methods("POST")
	serveMux.HandleFunc("/api/v1/deletetest", s.HandleApiDeleteTest).Methods("POST")

	serveMux.HandleFunc("/edit_test", s.HandleEditTest)
	serveMux.HandleFunc("/register", s.HandleRegister)
	serveMux.HandleFunc("/login", s.HandleLogin)
	serveMux.HandleFunc("/test", s.HandleGetTest)
	serveMux.HandleFunc("/", s.HandleMain)

	serveMux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("assets/static"))))
}
