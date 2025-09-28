package scraper

import (
	"context"
	"log"
	"sync"
	"time"

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
	log.Printf("Scraping feed: %v", feed.Url)
	_, err := s.Database.MarkFeedAsFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("Error marking feed as fetched: %v", err)
	}
	rssFeed, err := models.GetRSSFeedFromURL(feed.Url)
	if err != nil {
		log.Printf("Error getting RSS feed from URL: %v", err)
	}
	log.Printf("RSS feed: %v", rssFeed.Channel.Title)
	items := rssFeed.Channel.Items
	for _, item := range items {
		log.Printf("RSS feed %v found post: %v", feed.Url, item.Title)
	}
	log.Printf("Feed %v processed successfully, %v posts found", feed.Url, len(items))
}
