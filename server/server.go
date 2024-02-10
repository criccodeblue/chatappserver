package server

import (
	"chatappserver/database"
	"chatappserver/internal/model"
	"chatappserver/internal/util"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type Server struct {
	*mux.Router
	storage *database.PostgresStorage
	Port    string
}

func NewServer(port string, storage *database.PostgresStorage) *Server {
	s := &Server{
		Router:  mux.NewRouter(),
		storage: storage,
		Port:    port,
	}

	s.routes()

	return s
}

func (s *Server) routes() {
	s.HandleFunc("/users", s.GetUsers).Methods("GET")
	s.HandleFunc("/user/{handle}", s.VerifyAuth(s.GetUserByHandle)).Methods("GET")
	s.HandleFunc("/signup_user", s.SignUpUser).Methods("POST")
	s.HandleFunc("/login_user", s.LoginUser).Methods("POST")
	s.HandleFunc("/user", s.SignUpUser).Methods("POST")
	s.HandleFunc("/", s.Ping).Methods("GET")
}

func (s *Server) Ping(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode("Pong"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) VerifyAuth(endpointFunc func(http.ResponseWriter, *http.Request)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header["token"] != nil {
			fmt.Println(r.Header["token"][0])
			// validate token
			// if all ok
			endpointFunc(w, r)
			// else show error
		}
	}
}

func (s *Server) GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dbUsers, err := s.storage.GetUsers()

	if err != nil {
		util.ErrorResponse(w, err.Error())
	}

	if err := json.NewEncoder(w).Encode(dbUsers); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetUserByHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	handle, ok := vars["handle"]

	if !ok {
		http.Error(w, "Parsing Error", http.StatusInternalServerError)
		return
	}

	user, err := s.storage.GetUserByHandle(handle)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) SignUpUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var authUser model.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	hashedPassword, err := util.GetPasswordHash(authUser.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	authUser.PasswordHash = hashedPassword

	user, err := s.storage.CreateUser(&authUser)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var authUser model.AuthUser

	if err := json.NewDecoder(r.Body).Decode(&authUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	hashedPassword, err := s.storage.GetUserPassword(authUser.Handle)
	if err != nil {
		util.ErrorResponse(w, err.Error())
		return
	}

	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(authUser.Password)); err != nil {
		util.ErrorResponse(w, "Invalid Password")
		return
	} else {
		user, err := s.storage.GetUserByHandle(authUser.Handle)
		if err != nil {
			util.ErrorResponse(w, err.Error())
		}
		authToken, err := util.CreateToken(user.ID)
		if err != nil {
			util.ErrorResponse(w, err.Error())
		}
		response := model.ApiResponse{
			Status:  model.StatusOk,
			Message: "",
			Data:    authToken,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
