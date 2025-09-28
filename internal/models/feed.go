package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

type CreateFeedParams struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func (b *CreateFeedParams) Decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (b *CreateFeedParams) Validate() error {
	if b.Name == "" {
		return errors.New("name is required")
	}
	if b.Url == "" {
		return errors.New("url is required")
	}
	return nil
}

type Feed struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Url       string    `json:"url"`
	UserID    uuid.UUID `json:"user_id"`
}

func NewFeedFromDatabase(feed database.Feed) *Feed {
	f := &Feed{}
	f.ID = feed.ID
	f.CreatedAt = feed.CreatedAt
	f.UpdatedAt = feed.UpdatedAt
	f.Name = feed.Name
	f.Url = feed.Url
	f.UserID = feed.UserID
	return f
}
