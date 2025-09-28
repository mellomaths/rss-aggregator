package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

type CreateUserParams struct {
	Name string `json:"name"`
}

func (b *CreateUserParams) Decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (b *CreateUserParams) Validate() error {
	if b.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	ApiKey    string    `json:"api_key"`
}

func NewUserFromDatabase(user database.User) *User {
	u := &User{}
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.Name = user.Name
	u.ApiKey = user.ApiKey
	return u
}
