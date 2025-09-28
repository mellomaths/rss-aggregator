package models

import (
	"encoding/json"
	"encoding/xml"
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

// RSSFeed represents the basic structure of an RSS feed for validation
type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title       string `xml:"title"`
		Description string `xml:"description"`
		Link        string `xml:"link"`
	} `xml:"channel"`
}

func (b *RSSFeed) Decode(resp *http.Response) error {
	decoder := xml.NewDecoder(resp.Body)
	err := decoder.Decode(b)
	if err != nil {
		return err
	}
	return nil
}

func (b *RSSFeed) Validate() error {
	if b.XMLName.Local != "rss" {
		return errors.New("not a valid RSS feed")
	}
	return nil
}

// validateRSSURL checks if the provided URL points to a valid RSS feed
func validateRSSURL(url string) error {
	client := &http.Client{
		Timeout: 200 * time.Millisecond,
	}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("URL returned status code: %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "xml") &&
		!strings.Contains(strings.ToLower(contentType), "rss") &&
		!strings.Contains(strings.ToLower(contentType), "atom") {
		return fmt.Errorf("URL does not appear to be an RSS feed (content-type: %s)", contentType)
	}
	rssFeed := RSSFeed{}
	if err := rssFeed.Decode(resp); err != nil {
		return fmt.Errorf("failed to parse as RSS feed: %v", err)
	}
	if err := rssFeed.Validate(); err != nil {
		return fmt.Errorf("failed to validate RSS feed: %v", err)
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
	// Validate that the URL points to a valid RSS feed
	if err := validateRSSURL(b.Url); err != nil {
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
