package models

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

type CreateFeedFollowParams struct {
	FeedID uuid.UUID `json:"feed_id"`
}

func (b *CreateFeedFollowParams) Decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(b); err != nil {
		return err
	}
	return nil
}

type DeleteFeedFollowParams struct {
	ID uuid.UUID `json:"id"`
}

func (b *DeleteFeedFollowParams) Decode(feedFollowID string) error {
	if feedFollowID == "" {
		return errors.New("id is required")
	}
	id, err := uuid.Parse(feedFollowID)
	if err != nil {
		return err
	}
	b.ID = id
	return nil
}

type FeedFollow struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func NewFeedFollowFromDatabase(feedFollow database.FeedFollow) *FeedFollow {
	return &FeedFollow{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	}
}

func NewFeedFollowsFromDatabase(feedFollows []database.FeedFollow) []*FeedFollow {
	fs := make([]*FeedFollow, len(feedFollows))
	for i, feedFollow := range feedFollows {
		fs[i] = NewFeedFollowFromDatabase(feedFollow)
	}
	return fs
}
