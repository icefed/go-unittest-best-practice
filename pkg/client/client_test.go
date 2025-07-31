package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleFunc(w, r)
	}))
	defer server.Close()

	c := New(server.URL)

	t.Run("create", func(t *testing.T) {
		handleFunc = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`))
		}
		user, err := c.UserCreate(User{
			Name:  "liuliu",
			Email: "aa@bb.com",
		})
		assert.Nil(t, err)
		assert.EqualValues(t, &User{
			ID:        "0198271f-bc9d-74ac-a63b-41cf2c6c2f82",
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: time.Unix(1752999201, 0),
			UpdatedAt: time.Unix(1752999201, 0),
		}, user)
	})
	t.Run("get", func(t *testing.T) {
		handleFunc = func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{"data":{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}}`))
		}
		user, err := c.UserGet("0198271f-bc9d-74ac-a63b-41cf2c6c2f82")
		assert.Nil(t, err)
		assert.EqualValues(t, &User{
			ID:        "0198271f-bc9d-74ac-a63b-41cf2c6c2f82",
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: time.Unix(1752999201, 0),
			UpdatedAt: time.Unix(1752999201, 0),
		}, user)
	})
	t.Run("update", func(t *testing.T) {
		handleFunc = func(w http.ResponseWriter, r *http.Request) {
		}
		err := c.UserUpdate(User{
			ID:   "0198271f-bc9d-74ac-a63b-41cf2c6c2f82",
			Name: "liuliu2",
		})
		assert.Nil(t, err)
	})
	t.Run("delete", func(t *testing.T) {
		handleFunc = func(w http.ResponseWriter, r *http.Request) {
		}
		err := c.UserDelete("0198271f-bc9d-74ac-a63b-41cf2c6c2f82")
		assert.Nil(t, err)
	})
	t.Run("list", func(t *testing.T) {
		patches := gomonkey.NewPatches()
		defer patches.Reset()
		patches.ApplyMethodReturn(&http.Client{}, "Do", &http.Response{
			StatusCode: http.StatusOK,
		}, nil)
		patches.ApplyFuncReturn(io.ReadAll, []byte(`{"data":{"total":1,"users":[{"id":"0198271f-bc9d-74ac-a63b-41cf2c6c2f82","name":"liuliu","email":"aa@bb.com","age":0,"createdAt":"2025-07-20T16:13:21+08:00","updatedAt":"2025-07-20T16:13:21+08:00"}]}}`), nil)

		users, total, err := c.UserList()
		assert.Nil(t, err)
		assert.EqualValues(t, 1, total)
		assert.EqualValues(t, []User{{
			ID:        "0198271f-bc9d-74ac-a63b-41cf2c6c2f82",
			Name:      "liuliu",
			Email:     "aa@bb.com",
			CreatedAt: time.Unix(1752999201, 0),
			UpdatedAt: time.Unix(1752999201, 0),
		}}, users)
	})
}