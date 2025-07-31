package fake

import (
	"fmt"
	"go-unittest-best-practice/pkg/client"
	"sync"
	"time"

	"github.com/google/uuid"
)

type fakeClient struct {
	mu           sync.Mutex
	users        map[string]*client.User
	usersByEmail map[string]string
}

var _ client.Client = &fakeClient{}

func (c *fakeClient) UserCreate(u client.User) (*client.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.usersByEmail[u.Email]
	if ok {
		return nil, fmt.Errorf("user already exists")
	}

	id := uuid.NewString()
	now := time.Now()
	c.users[id] = &client.User{
		ID:        id,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		CreatedAt: now,
		UpdatedAt: now,
	}
	c.usersByEmail[u.Email] = id
	return c.users[id], nil
}

func (c *fakeClient) UserGet(id string) (*client.User, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	user, ok := c.users[id]
	if !ok {
		return nil, fmt.Errorf("user already exists")
	}
	return user, nil
}

func (c *fakeClient) UserUpdate(u client.User) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	user, ok := c.users[u.ID]
	if !ok {
		return fmt.Errorf("user already exists")
	}
	user.Name = u.Name
	return nil
}

func (c *fakeClient) UserDelete(id string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	user, ok := c.users[id]
	if !ok {
		return fmt.Errorf("user already exists")
	}
	delete(c.usersByEmail, user.Email)
	delete(c.users, id)
	return nil
}

func (c *fakeClient) UserList() ([]client.User, int64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	total := len(c.users)
	users := make([]client.User, 0, total)
	for _, u := range c.users {
		users = append(users, client.User{
			ID:        u.ID,
			Name:      u.Name,
			Email:     u.Email,
			Age:       u.Age,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	return users, int64(total), nil
}