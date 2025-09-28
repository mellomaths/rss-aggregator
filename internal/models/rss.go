package models

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RSSFeed struct {
	XMLName xml.Name `xml:"rss"`
	Version string   `xml:"version,attr"`
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Language    string    `xml:"language"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (b *RSSFeed) Validate() error {
	if b.XMLName.Local != "rss" {
		return errors.New("not a valid RSS feed")
	}
	return nil
}

func GetRSSFeedFromURL(url string) (RSSFeed, error) {
	client := &http.Client{
		Timeout: 200 * time.Millisecond,
	}
	resp, err := client.Get(url)
	if err != nil {
		return RSSFeed{}, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return RSSFeed{}, fmt.Errorf("URL returned status code: %d", resp.StatusCode)
	}
	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(strings.ToLower(contentType), "xml") &&
		!strings.Contains(strings.ToLower(contentType), "rss") &&
		!strings.Contains(strings.ToLower(contentType), "atom") {
		return RSSFeed{}, fmt.Errorf("URL does not appear to be an RSS feed (content-type: %s)", contentType)
	}
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return RSSFeed{}, fmt.Errorf("failed to read RSS feed: %v", err)
	}
	rssFeed := RSSFeed{}
	if err := xml.Unmarshal(dat, &rssFeed); err != nil {
		return RSSFeed{}, fmt.Errorf("failed to unmarshal XML: %v", err)
	}
	return rssFeed, nil
}
