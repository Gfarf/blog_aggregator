package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"slices"
	"time"

	"github.com/Gfarf/blog_aggregator/internal/database"
	"github.com/google/uuid"
)

// RSS Feed XML struct
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

// Client -
type Client struct {
	httpClient http.Client
}

// NewClient -
func NewClient(timeout time.Duration) Client {
	return Client{
		httpClient: http.Client{
			Timeout: timeout,
		},
	}
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	c := NewClient(20 * time.Second)
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rssFeed := RSSFeed{}
	err = xml.Unmarshal(dat, &rssFeed)
	if err != nil {
		return nil, err
	}
	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i, item := range rssFeed.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Description = html.UnescapeString(item.Description)
		rssFeed.Channel.Item[i] = item
	}
	return &rssFeed, nil
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return (func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	})
}

func scrapeFeeds(s *state) error {
	//Get the next feed
	feedID, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}
	err = s.db.MarkFeedFetched(
		context.Background(),
		database.MarkFeedFetchedParams{ID: feedID, UpdatedAt: time.Now(),
			LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}})
	if err != nil {
		return err
	}
	feed, err := s.db.GetFeedByID(context.Background(), feedID)
	if err != nil {
		return err
	}
	fmt.Println(feed.Name)
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		return err
	}
	postsUrl, err := s.db.GetPostsUrls(context.Background())
	if err != nil {
		return err
	}
	for _, item := range feedData.Channel.Item {
		if slices.Contains(postsUrl, item.Link) {
			continue
		}
		parsedTime, err := timeParser(item.PubDate)
		if err != nil {
			fmt.Println(err)
		}
		_, err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       sql.NullString{String: item.Title, Valid: true},
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: parsedTime,
			FeedID:      feedID,
		})
		if err != nil {
			return fmt.Errorf("fail during Post creation: %v", err)
		}
	}
	return nil
}

func timeParser(s string) (time.Time, error) {
	res, err := time.Parse(time.RFC3339, s)
	if err != nil {
		res, err = time.Parse(time.ANSIC, s)
		if err != nil {
			res, err = time.Parse(time.UnixDate, s)
			if err != nil {
				res, err = time.Parse("Sat, 25 Jul 2020 00:00:00 +0000", s)
				if err != nil {
					var t time.Time
					return t, fmt.Errorf("error parsing time: %v", err)
				}
				return res, nil
			}
			return res, nil
		}
		return res, nil
	}
	return res, nil
}
