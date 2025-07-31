package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"time"

	"go.uber.org/mock/gomock"

	"go-unittest-best-practice/internal/store"
)

func (s *ServiceTestSuite) TestCreateUser() {
	s.Run("param name empty", func() {
		req := httptest.NewRequest("POST", "http://127.0.0.1:8888/user/create?email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusBadRequest, w.Code)
		s.EqualValues(`{"error":"param name not set"}`, w.Body.String())
	})
	s.Run("success", func() {
		t := time.Unix(1752999201, 0)
		id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
		s.mockUserRepo.EXPECT().Create(gomock.Any()).Return(nil).Times(1)
		s.mockUserRepo.EXPECT().GetByEmail(gomock.Any()).Return(&store.User{
			ID:        id,
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: t,
			UpdatedAt: t,
		}, nil).Times(1)

		req := httptest.NewRequest("POST", "http://127.0.0.1:8888/user/create?name=liuliu&email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusOK, w.Code)
		s.EqualValues(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`, w.Body.String())
	})
}

func (s *ServiceTestSuite) TestGetUser() {
	s.Run("user not found", func() {
		s.mockUserRepo.EXPECT().GetByEmail(gomock.Any()).Return(nil, fmt.Errorf("user not found")).Times(1)
		req := httptest.NewRequest("GET", "http://127.0.0.1:8888/user/get?email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusInternalServerError, w.Code)
		s.EqualValues(`{"error":"user not found"}`, w.Body.String())
	})
	s.Run("success", func() {
		t := time.Unix(1752999201, 0)
		id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
		s.mockUserRepo.EXPECT().GetByEmail(gomock.Any()).Return(&store.User{
			ID:        id,
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: t,
			UpdatedAt: t,
		}, nil).Times(1)

		req := httptest.NewRequest("GET", "http://127.0.0.1:8888/user/get?email=aa@bb.com", nil)
		w := httptest.NewRecorder()
		s.svc.ServeHTTP(w, req)
		s.EqualValues(http.StatusOK, w.Code)
		s.EqualValues(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`, w.Body.String())
	})
}

func (s *ServiceTestSuite) TestUpdateUser() {
	t := time.Unix(1752999201, 0)
	id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
	s.mockUserRepo.EXPECT().GetByID(id).Return(&store.User{
		ID:        id,
		Name:      "liuliu",
		Email:     "aa@bb.com",
		CreatedAt: t,
		UpdatedAt: t,
	}, nil).Times(1)
	s.mockUserRepo.EXPECT().Update(gomock.Any()).Return(nil).Times(1)

	req := httptest.NewRequest("POST", "http://127.0.0.1:8888/user/update?id=0198271f-bc9d-74ac-a63b-41cf2c6c2f82&name=liuliu2", nil)
	w := httptest.NewRecorder()
	s.svc.ServeHTTP(w, req)
	s.EqualValues(http.StatusOK, w.Code)
	s.EqualValues(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu2","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`, w.Body.String())
}

func (s *ServiceTestSuite) TestDeleteUser() {
	id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
	s.mockUserRepo.EXPECT().DeleteByID(id).Return(nil).Times(1)

	req := httptest.NewRequest("POST", "http://127.0.0.1:8888/user/delete?id=0198271f-bc9d-74ac-a63b-41cf2c6c2f82", nil)
	w := httptest.NewRecorder()
	s.svc.ServeHTTP(w, req)
	s.EqualValues(http.StatusOK, w.Code)
}

func (s *ServiceTestSuite) TestListUser() {
	t := time.Unix(1752999201, 0)
	id := "0198271f-bc9d-74ac-a63b-41cf2c6c2f82"
	s.mockUserRepo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]store.User{{
		ID:        id,
		Name:      "liuliu",
		Email:     "aa@bb.com",
		CreatedAt: t,
		UpdatedAt: t,
	}}, int64(1), nil).Times(1)

	req := httptest.NewRequest("GET", "http://127.0.0.1:8888/user/list", nil)
	w := httptest.NewRecorder()
	s.svc.ServeHTTP(w, req)
	s.EqualValues(http.StatusOK, w.Code)
	s.EqualValues(`{"data":{"total":1,"users":[{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}]}}`, w.Body.String())
}