package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Url         string    `json:"url"`
	Description string    `json:"description"`
	PublishedAt time.Time `json:"published_at"`
	FeedID      uuid.UUID `json:"feed_id"`
}

func NewPostFromDatabase(post database.Post) *Post {
	return &Post{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Url:         post.Url,
		Description: post.Description.String,
		PublishedAt: post.PublishedAt,
		FeedID:      post.FeedID,
	}
}

func NewPostsFromDatabase(posts []database.Post) []*Post {
	ps := make([]*Post, len(posts))
	for i, post := range posts {
		ps[i] = NewPostFromDatabase(post)
	}
	return ps
}
