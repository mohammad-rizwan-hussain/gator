package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/mohammad-rizwan-hussain/gator/internal/config"
	"github.com/mohammad-rizwan-hussain/gator/internal/database"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Get the feed

	client := &http.Client{}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, feedURL, nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error framing the request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error requesting from url: %s\nerror: %w", feedURL, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error reading resp body: %w", err)
	}
	feed := RSSFeed{}
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error unmarshaling resp body: %w", err)
	}

	// decode escaped HTML
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i := range feed.Channel.Item {
		feed.Channel.Item[i].Title = html.UnescapeString(feed.Channel.Item[i].Title)
		feed.Channel.Item[i].Description = html.UnescapeString(feed.Channel.Item[i].Description)
	}

	return &feed, nil
}

func ScrapeFeeds(s *config.State) error {
	nextFeed, err := s.DB.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("\nFetching feeds from %s...\n", nextFeed.Name)

	feed, err := FetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		return err
	}

	_, err = s.DB.MarkFeedFetched(context.Background(), nextFeed.ID)
	if err != nil {
		return err
	}

	for _, feedItem := range feed.Channel.Item {
		// fmt.Printf("- %s\n", feedItem.Title)
		publishedAt, err := time.Parse(time.RFC1123Z, feedItem.PubDate)
		if err != nil {
			return err
		}

		post := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       feedItem.Title,
			Url:         feedItem.Link,
			Description: feedItem.Description,
			PublishedAt: publishedAt,
			FeedID:      nextFeed.ID,
		}
		_, err = s.DB.CreatePost(context.Background(), post)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Post Created!")
	}

	return nil
}
