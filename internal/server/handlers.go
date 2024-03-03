package server

import (
	"fmt"
	"html/template"
	"net/http"
)

func (s *Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
	s.logger.Info("handled /api/ping")
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/registration.html")
	tmpl.Execute(w, nil)
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/authorization.html")
	tmpl.Execute(w, nil)
}

func (s *Server) HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "There's all routes for now that handles and works:\n- /login\n- /register")
}
