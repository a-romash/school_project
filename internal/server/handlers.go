package server

import (
	"encoding/json"
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
	s.logger.Info("handled /register")
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/authorization.html")
	tmpl.Execute(w, nil)
	s.logger.Info("handled /login")
}

func (s *Server) HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "There's all routes for now that handles and works:\n- /login\n- /register")
	s.logger.Info("handled /")
}

func (s *Server) HandleAPILogin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Login string
		Pw    string
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		s.logger.Error("Error on /api/auth:\n" + err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
	s.logger.Info("handled /api/auth")
}
