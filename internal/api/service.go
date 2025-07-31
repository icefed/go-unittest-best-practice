package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go-unittest-best-practice/internal/config"
	"go-unittest-best-practice/internal/store"

	"github.com/google/uuid"
)

type Service struct {
	mux  *http.ServeMux
	conf *config.Config

	userRepo store.UserRepository
}

func NewService(userRepo store.UserRepository, conf *config.Config) *Service {
	mux := http.NewServeMux()
	service := &Service{
		mux:      mux,
		conf:     conf,
		userRepo: userRepo,
	}
	mux.HandleFunc("/user/create", service.createUser)
	mux.HandleFunc("/user/get", service.getUser)
	mux.HandleFunc("/user/update", service.updateUser)
	mux.HandleFunc("/user/delete", service.deleteUser)
	mux.HandleFunc("/user/list", service.listUser)
	return service
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Service) createUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.error(w, fmt.Errorf("param name not set"))
		return
	}
	email := r.FormValue("email")
	if email == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.error(w, fmt.Errorf("param email not set"))
		return
	}

	err := s.userRepo.Create(&store.User{
		ID:    uuid.NewString(),
		Name:  name,
		Email: email,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
		return
	}

	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
		return
	}
	s.data(w, convertModelUser(user))
}

func (s *Service) getUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	email := r.FormValue("email")
	if id == "" && email == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.error(w, fmt.Errorf("param id or email not set"))
		return
	}
	var user *store.User
	var err error
	if id != "" {
		user, err = s.userRepo.GetByID(id)
	} else {
		user, err = s.userRepo.GetByEmail(email)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
		return
	}
	s.data(w, convertModelUser(user))
}

func (s *Service) updateUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.error(w, fmt.Errorf("param id not set"))
		return
	}
	name := r.FormValue("name")
	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.error(w, fmt.Errorf("param name not set"))
		return
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
		return
	}

	user.Name = name
	err = s.userRepo.Update(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
		return
	}

	s.data(w, convertModelUser(user))
}

func (s *Service) deleteUser(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		s.error(w, fmt.Errorf("param id not set"))
		return
	}
	err := s.userRepo.DeleteByID(id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
	}
}

func (s *Service) listUser(w http.ResponseWriter, r *http.Request) {
	users, total, err := s.userRepo.List(1, 100)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		s.error(w, err)
		return
	}
	dataUsers := make([]User, 0, len(users))
	for _, u := range users {
		dataUsers = append(dataUsers, *convertModelUser(&u))
	}
	s.data(w, map[string]interface{}{
		"total": total,
		"users": dataUsers,
	})
}

func (s *Service) data(w http.ResponseWriter, body interface{}) {
	data, _ := json.Marshal(&DataResponse{Data: body})
	w.Write(data)
}

func (s *Service) error(w http.ResponseWriter, err error) {
	data, _ := json.Marshal(&ErrorResponse{Error: err.Error()})
	w.Write(data)
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type DataResponse struct {
	Data interface{} `json:"data"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func convertModelUser(modelUser *store.User) *User {
	return &User{
		ID:        modelUser.ID,
		Name:      modelUser.Name,
		Email:     modelUser.Email,
		Age:       modelUser.Age,
		CreatedAt: modelUser.CreatedAt,
		UpdatedAt: modelUser.UpdatedAt,
	}
}
