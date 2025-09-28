package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	isValidUrl := strings.HasPrefix(b.Url, "https://") || strings.HasPrefix(b.Url, "http://")
	if !isValidUrl {
		return errors.New("url must start with https:// or http://")
	}
	_, err := GetRSSFeedFromURL(b.Url)
	if err != nil {
		return fmt.Errorf("invalid RSS URL: %v", err)
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
	return &Feed{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.Url,
		UserID:    feed.UserID,
	}
}

func NewFeedsFromDatabase(feeds []database.Feed) []*Feed {
	fs := make([]*Feed, len(feeds))
	for i, feed := range feeds {
		fs[i] = NewFeedFromDatabase(feed)
	}
	return fs
}
