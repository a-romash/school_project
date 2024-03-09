package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"project/pkg/database/postgresql"
	"project/pkg/model"
	"time"

	"github.com/google/uuid"
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
	string_token := r.Header.Get("t")
	slog.Debug("token: " + string_token)
	if string_token != "null" && string_token != "" {
		session_token, _ := uuid.Parse(string_token)
		err := s.db.ValidateToken(session_token)
		if err != nil {
			if _, ok := err.(postgresql.ErrTokenExpired); ok {
				slog.Debug("token is expired: " + session_token.String())
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			slog.Error("error while validating token")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("Error on /api/login:\n" + err.Error())
			return
		}
		slog.Debug("token is ok")
		w.WriteHeader(http.StatusOK)
		return
	}

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

	sessionToken, err := s.db.NewToken(user.Login)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		slog.Error("Error on /api/login:\n" + err.Error())
		return
	}

	var t struct {
		Login      string    `json:"login"`
		Token      string    `json:"t"`
		Expires_at time.Time `json:"expires_at"`
	}

	t.Login = sessionToken.Login
	t.Token = sessionToken.Token.String()
	t.Expires_at = sessionToken.Expires_at

	byteToken, err := json.Marshal(t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/login:\n" + err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(byteToken)
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
