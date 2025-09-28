package scraper

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mellomaths/rss-aggregator/internal/database"
	"github.com/mellomaths/rss-aggregator/internal/models"
)

type RSSScraper struct {
	Database            *database.Queries
	Concurrency         int
	TimeBetweenRequests time.Duration
}

func (s *RSSScraper) Start() {
	log.Printf("Scraping %v concurrent requests every %v duration", s.Concurrency, s.TimeBetweenRequests)
	ticker := time.NewTicker(s.TimeBetweenRequests)
	for ; ; <-ticker.C {
		feeds, err := s.Database.GetNextFeedsToFetch(context.Background(), int32(s.Concurrency))
		if err != nil {
			log.Printf("Error getting next feeds to fetch: %v", err)
			continue
		}

		wg := sync.WaitGroup{}
		for _, feed := range feeds {
			wg.Add(1)
			go s.scrapeFeed(&wg, &feed)
		}
		wg.Wait()
	}
}

func (s *RSSScraper) scrapeFeed(wg *sync.WaitGroup, feed *database.Feed) {
	defer wg.Done()
	log.Printf("Scraping feed %v (%v)", feed.Name, feed.ID)
	rssFeed, err := models.GetRSSFeedFromURL(feed.Url)
	if err != nil {
		log.Printf("Error getting RSS feed from URL %v: %v", feed.Url, err)
	}
	log.Printf("Processing RSS feed %v (%v)", rssFeed.Channel.Title, feed.ID)
	items := rssFeed.Channel.Items
	for _, item := range items {
		description := sql.NullString{}
		if item.Description != "" {
			description.String = item.Description
			description.Valid = true
		}
		publishedAt, err := time.Parse(time.RFC1123Z, item.PubDate)
		if err != nil {
			log.Printf("Error parsing published date %v: %v", item.PubDate, err)
			continue
		}
		_, err = s.Database.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: publishedAt,
			FeedID:      feed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				continue
			}
			log.Printf("Error creating post %v (%v): %v", item.Title, item.Link, err)
			continue
		}
	}
	log.Printf("Feed %v (%v) processed successfully, %v posts found", feed.Name, feed.ID, len(items))
	_, err = s.Database.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched %v (%v): %v", feed.Name, feed.ID, err)
	}
}
