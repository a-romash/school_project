package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"project/pkg/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Server) HandlePing(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("ok"))
	slog.Info("handled /api/ping")
}

func (s *Server) HandleRegister(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/registration.html")
	tmpl.Execute(w, nil)
	slog.Info("handled /register")
}

func (s *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/authorization.html")
	tmpl.Execute(w, nil)
	slog.Info("handled /login")
}

func (s *Server) HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "There's all routes for now that handles and works:\n- /login\n- /register")
	slog.Info("handled /")
}

func (s *Server) HandleAPILogin(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Login string `json:"login"`
		Pw    string `json:"pw"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/login:\n" + err.Error())
		return
	}

	user, err := s.db.GetUser(data.Login)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/login:\n" + err.Error())
		return
	}
	if !user.CheckPassword(data.Pw) {
		http.Error(w, errors.New("password isn't correct").Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("handled /api/login")
}

func (s *Server) HandleAPIRegister(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Login    string `json:"login"`
		Name     string `json:"name"`
		Lastname string `json:"lastname"`
		School   string `json:"school"`
		Pw       string `json:"pw"`
	}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/register:\n" + err.Error())
		return
	}

	newUser, err := model.CreateUser(data.Login, data.Name, data.Lastname, data.School, data.Pw)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/register:\n" + err.Error())
		return
	}

	err = s.db.RegisterNewUser(newUser)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/register:\n" + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	slog.Info("handled /api/register")
	slog.Info("registered new user:\nName: " + newUser.Name + "\nLastname: " + newUser.Lastname + "\nSchool: " + newUser.School + "\nLogin: " + newUser.Login)
}
