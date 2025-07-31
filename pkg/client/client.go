package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

//go:generate mockgen -source=client.go -destination=client_mock.go -package=client
type Client interface {
	UserCreate(u User) (*User, error)
	UserGet(id string) (*User, error)
	UserUpdate(u User) error
	UserDelete(id string) error
	UserList() ([]User, int64, error)
}

var _ Client = &client{}

func New(server string) Client {
	c := &client{
		httpClient: &http.Client{},
		server:     server,
	}
	return c
}

type client struct {
	httpClient *http.Client
	server     string
}

func (c *client) UserCreate(u User) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/create", c.server), nil)
	if err != nil {
		return nil, err
	}
	postData := url.Values{}
	postData.Set("name", u.Name)
	postData.Set("email", u.Email)
	req.Form = postData

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %s", string(data))
	}
	if resp.StatusCode != http.StatusOK {
		var errRes Response
		json.Unmarshal(data, &errRes)
		if errRes.Error != "" {
			return nil, errors.New(errRes.Error)
		}
		return nil, fmt.Errorf("create user failed: %d", resp.StatusCode)
	}

	var userResp CreateGetResponse
	err = json.Unmarshal(data, &userResp)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %v", err)
	}
	return &userResp.Data, nil
}

func (c *client) UserGet(id string) (*User, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/create", c.server), nil)
	if err != nil {
		return nil, err
	}
	postData := url.Values{}
	postData.Set("id", id)
	req.Form = postData

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("create user failed: %s", string(data))
	}
	if resp.StatusCode != http.StatusOK {
		var errRes Response
		json.Unmarshal(data, &errRes)
		if errRes.Error != "" {
			return nil, errors.New(errRes.Error)
		}
		return nil, fmt.Errorf("create user failed: %d", resp.StatusCode)
	}

	var userResp CreateGetResponse
	err = json.Unmarshal(data, &userResp)
	if err != nil {
		return nil, fmt.Errorf("parse response failed: %v", err)
	}
	return &userResp.Data, nil
}

func (c *client) UserUpdate(u User) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/update", c.server), nil)
	if err != nil {
		return nil
	}
	postData := url.Values{}
	postData.Set("name", u.Name)
	postData.Set("id", u.ID)
	req.Form = postData

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errRes Response
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("update user failed: %s", string(data))
		}
		json.Unmarshal(data, &errRes)
		if errRes.Error != "" {
			return errors.New(errRes.Error)
		}
		return fmt.Errorf("update user failed: %d", resp.StatusCode)
	}

	return nil
}

func (c *client) UserDelete(id string) error {
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/user/delete", c.server), nil)
	if err != nil {
		return err
	}
	postData := url.Values{}
	postData.Set("id", id)
	req.Form = postData

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errRes Response
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("delete user failed: %s", string(data))
		}
		json.Unmarshal(data, &errRes)
		if errRes.Error != "" {
			return errors.New(errRes.Error)
		}
		return fmt.Errorf("delete user failed: %d", resp.StatusCode)
	}

	return nil
}

func (c *client) UserList() ([]User, int64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/user/list", c.server), nil)
	if err != nil {
		return nil, 0, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("list user failed: %s", string(data))
	}
	if resp.StatusCode != http.StatusOK {
		var errRes Response
		json.Unmarshal(data, &errRes)
		if errRes.Error != "" {
			return nil, 0, errors.New(errRes.Error)
		}
		return nil, 0, fmt.Errorf("list user failed: %d", resp.StatusCode)
	}

	var listResp ListResponse
	err = json.Unmarshal(data, &listResp)
	if err != nil {
		return nil, 0, fmt.Errorf("parse response failed: %v", err)
	}
	return listResp.Data.Users, listResp.Data.Total, nil
}

type CreateGetResponse struct {
	Data User `json:"data"`
}

type ListResponse struct {
	Data ListResponseData `json:"data"`
}

type ListResponseData struct {
	Total int64  `json:"total"`
	Users []User `json:"users"`
}

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
