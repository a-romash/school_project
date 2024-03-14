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

func (s *Server) HandleGetTest(w http.ResponseWriter, r *http.Request) {
	testId := r.URL.Query().Get("t")
	if testId == "" {
		http.Error(w, "parameter 't' is required", http.StatusBadRequest)
		return
	}

	tmpl, _ := template.ParseFiles("assets/templates/test.html")
	tmpl.Execute(w, nil)
	slog.Info("handled /test?t=" + testId)
}

func (s *Server) HandleMain(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/main.html")
	tmpl.Execute(w, nil)
	slog.Info("handled /")
}

type ErrTokenUndefined error

func (s *Server) ParseToken(string_token string) (token uuid.UUID, err error) {
	if string_token != "null" && string_token != "" {
		token, _ = uuid.Parse(string_token)
		err = s.db.ValidateToken(token)
		if err != nil {
			if _, ok := err.(postgresql.ErrTokenExpired); ok {
				slog.Debug("token is expired: " + token.String())
				return uuid.Nil, err
			}
			slog.Error("error while validating token")
			return uuid.Nil, err
		}
		slog.Debug("token is ok")
		return token, nil
	}
	var e ErrTokenUndefined = errors.New("token undefined")
	return uuid.Nil, e
}

func (s *Server) HandleAPILogin(w http.ResponseWriter, r *http.Request) {
	string_token := r.Header.Get("t")
	slog.Debug("token: " + string_token)
	_, err := s.ParseToken(string_token)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	} else if _, ok := err.(ErrTokenUndefined); !ok {
		switch err.(type) {
		case postgresql.ErrTokenExpired:
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		default:
			http.Error(w, err.Error(), http.StatusInternalServerError)
			slog.Error("Error on /api/login:\n" + err.Error())
			return
		}
	}

	var data struct {
		Login string `json:"login"`
		Pw    string `json:"pw"`
	}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		slog.Error("Error on /api/login:\n" + err.Error())
		return
	}

	user, err := s.db.GetUserByLogin(data.Login)
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

func (s *Server) HandleApiCreateNewTest(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var test *model.Test
	var userName string
	var userId int

	if t, ok := data["token"]; ok {
		token, err := s.ParseToken(t.(string))
		if err != nil {
			http.Error(w, "token is invalid", http.StatusUnauthorized)
			return
		}
		user, err := s.db.GetUserByToken(token)
		if err != nil {
			http.Error(w, "token is undefined", http.StatusUnauthorized)
			return
		}
		userName = fmt.Sprint(user.Name, " ", user.Lastname)
		userId = user.Id
	} else if login, ok := data["login"]; ok {
		user, err := s.db.GetUserByLogin(login.(string))
		if err != nil {
			http.Error(w, "login is undefined", http.StatusUnauthorized)
			return
		}
		userName = fmt.Sprint(user.Name, " ", user.Lastname)
		userId = user.Id
	} else {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	test = model.CreateTest(userName, userId, data["questions"].([]interface{}), data["answers"].([]interface{}))
	fmt.Println(test.Questions...)
	err = s.db.CreateNewTest(test)
	if err != nil {
		slog.Error("error while creating test")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) HandleApiGetTest(w http.ResponseWriter, r *http.Request) {
	testId := r.URL.Query().Get("t")
	if testId == "" {
		http.Error(w, "parameter 't' is required", http.StatusBadRequest)
		return
	}

	test, err := s.db.GetTest(testId)
	if err != nil {
		http.Error(w, "test isn't exist", http.StatusBadRequest)
		return
	}

	data, err := test.GetJson()
	if err != nil {
		return
	}
	w.Write(data)
}
