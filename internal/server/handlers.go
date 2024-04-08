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
	"reflect"
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
	tmpl, _ := template.ParseFiles("assets/templates/authentication.html")
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

func (s *Server) HandleEditTest(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("assets/templates/create_test.html")
	tmpl.Execute(w, nil)
	slog.Info("handled /edittest")
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

	test = model.CreateTest(userName, data["title"].(string), userId, data["questions"].([]interface{}), data["answers"].([]interface{}))
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

func (s *Server) HandleApiGetInfo(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userId int
	var name string
	if t, ok := data["token"]; ok {
		token, err := s.ParseToken(t.(string))
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "token is invalid", http.StatusUnauthorized)
			return
		}
		user, err := s.db.GetUserByToken(token)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "token is undefined", http.StatusUnauthorized)
			return
		}
		userId = user.Id
		name = fmt.Sprint(user.Name, " ", user.Lastname)
		testIds, err := s.db.GetIdTests(userId)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			slog.Error("err getting test ids\n" + err.Error())
			http.Error(w, "error with test", http.StatusInternalServerError)
		}

		solutions := make(map[string]map[string]interface{})
		for _, v := range testIds {
			sols, err := s.db.GetSolutionsForTest(v)
			if err != nil && !errors.Is(err, pgx.ErrNoRows) {
				slog.Error("err while geting solutions\n" + err.Error())
				http.Error(w, "error with solutions", http.StatusInternalServerError)
			}
			solutions[v] = make(map[string]interface{})
			solutions[v]["solutions"] = [][]byte{}
			for _, sol := range sols {
				byteSolution, err := json.Marshal(sol)
				if err != nil {
					slog.Error("err marshalling solution\n" + err.Error())
					http.Error(w, "error with solution", http.StatusInternalServerError)
				}
				solutions[v]["solutions"] = append(solutions[v]["solutions"].([][]byte), byteSolution)
			}
			test, _ := s.db.GetTest(v)
			solutions[v]["test_name"] = test.Title
			solutions[v]["created"] = test.Created
			solutions[v]["updated"] = test.Updated
			solutions[v]["amount"] = len(solutions[v]["solutions"].([][]byte))
			solutions[v]["max_score"] = test.Max_score
		}
		response := make(map[string]interface{})
		response["user_id"] = userId
		response["name"] = name
		response["tests"] = solutions

		byteResponse, err := json.Marshal(response)
		if err != nil {
			slog.Error("err marshalling response" + err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Write(byteResponse)
		return
	}
	http.Error(w, "token is undefined", http.StatusUnauthorized)
}

func (s *Server) HandleApiGetResult(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	test, err := s.db.GetTest(data["test_id"].(string))
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	max_score := len(test.Answers)
	var cur_score int
	for i := range data["answers"].([]interface{}) {
		if reflect.DeepEqual(test.Answers[i], data["answers"].([]interface{})[i].(string)) {
			cur_score++
		}
	}
	response := make(map[string]interface{})
	response["test_id"] = data["test_id"].(string)
	response["cur_score"] = cur_score
	response["max_score"] = max_score

	byteResponse, err := json.Marshal(response)
	if err != nil {
		slog.Error("err marshalling response" + err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var answers []string
	for _, v := range data["answers"].([]interface{}) {
		answers = append(answers, v.(string))
	}

	solution := &model.Solution{
		Author:  data["author"].(string),
		Class:   data["class"].(string),
		Answers: answers,
		TestId:  uuid.MustParse(data["test_id"].(string)),
		Result:  cur_score,
	}

	err = s.db.CreateNewSolution(solution)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(byteResponse)
}

func (s *Server) HandleApiDeleteToken(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if t, ok := data["token"]; ok {
		token, err := uuid.Parse(t.(string))
		if err != nil {
			http.Error(w, "invalid token", http.StatusBadRequest)
		}
		_ = s.db.DeleteToken(token)
		w.WriteHeader(http.StatusOK)
	}
	http.Error(w, "token is empty", http.StatusBadRequest)
}

func (s *Server) HandleApiDeleteTest(w http.ResponseWriter, r *http.Request) {
	var data map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		slog.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if t, ok_t := data["token"]; ok_t {
		token, err := uuid.Parse(t.(string))
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}
		err = s.db.ValidateToken(token)
		if err != nil {
			slog.Error(err.Error())
			http.Error(w, "token is invalid or expired", http.StatusUnauthorized)
			return
		}
		if id, ok_id := data["test_id"]; ok_id {
			err = s.db.DeleteTest(id.(string))
			if err != nil {
				slog.Error(err.Error())
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "no test_id", http.StatusBadRequest)
		return
	}
	http.Error(w, "no token", http.StatusUnauthorized)
}
