package main

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/AlexTLDR/gator/internal/database"
)

type RSSFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Version string     `xml:"version,attr"`
	Channel RSSChannel `xml:"channel"`
}

type RSSChannel struct {
	Title         string    `xml:"title"`
	Link          string    `xml:"link"`
	Description   string    `xml:"description"`
	Language      string    `xml:"language"`
	PubDate       string    `xml:"pubDate"`
	LastBuildDate string    `xml:"lastBuildDate"`
	Items         []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	GUID        string   `xml:"guid"`
	Categories  []string `xml:"category"`
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	
	fmt.Printf("🌐 Fetching URL: %s\n", feedURL)

	req.Header.Set("User-Agent", "gator")

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("error parsing feed: %w", err)
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)

	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
		feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
	}

	return &feed, nil
}

func scrapeFeeds(s *state) error {
	ctx := context.Background()

	// Get the next feed to fetch
	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("no feeds available - add some feeds using the 'addfeed' command")
		}
		return fmt.Errorf("error getting next feed to fetch: %w", err)
	}

	// Mark the feed as fetched
	now := time.Now().UTC()
	err = s.db.MarkFeedFetched(ctx, database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		ID:            feed.ID,
	})
	if err != nil {
		return fmt.Errorf("error marking feed as fetched: %w", err)
	}

	// Fetch the feed content
	fmt.Printf("\n===== [%s] =====\n", now.Format(time.RFC3339))
	fmt.Printf("📥 Fetching feed: %s\n", feed.Name)
	fmt.Printf("🔗 URL: %s\n", feed.Url)
	
	if feed.LastFetchedAt.Valid {
		fmt.Printf("🕒 Last fetched: %s (%.1f hours ago)\n", 
			feed.LastFetchedAt.Time.Format(time.RFC3339),
			now.Sub(feed.LastFetchedAt.Time).Hours())
	} else {
		fmt.Printf("🕒 Last fetched: Never\n")
	}
	
	rssFeed, err := fetchFeed(ctx, feed.Url)
	if err != nil {
		return fmt.Errorf("error fetching feed content: %w", err)
	}

	// Print feed metadata
	fmt.Printf("📰 Title: %s\n", rssFeed.Channel.Title)
	if rssFeed.Channel.Description != "" {
		fmt.Printf("📝 Description: %s\n", rssFeed.Channel.Description)
	}
	
	// Print feed items
	fmt.Printf("📚 Found %d items in feed\n", len(rssFeed.Channel.Items))
	fmt.Println("----------------------------")
	for i, item := range rssFeed.Channel.Items {
		pubDate := item.PubDate
		if pubDate == "" {
			pubDate = "No date"
		}
		fmt.Printf("%d. [%s] %s\n", i+1, pubDate, item.Title)
		fmt.Printf("   🔗 %s\n", item.Link)
		if len(item.Categories) > 0 {
			fmt.Printf("   🏷️  %s\n", strings.Join(item.Categories, ", "))
		}
		fmt.Println()
	}
	fmt.Println("===========================")

	return nil
}

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <time_between_reqs>", cmd.Name)
	}

	// Parse the time duration
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format: %w (examples: 30s, 1m, 5m, 1h)", err)
	}

	fmt.Printf("🚀 Starting feed aggregation\n")
	fmt.Printf("⏱️  Collecting feeds every %s\n", timeBetweenRequests)
	fmt.Printf("📊 Feed collection started at: %s\n", time.Now().Format(time.RFC3339))
	fmt.Printf("❗ Press Ctrl+C to stop\n\n")

	// Create a ticker that triggers every timeBetweenRequests
	ticker := time.NewTicker(timeBetweenRequests)
	defer ticker.Stop()
	
	count := 0
	// Run immediately and then on each tick
	for {
		count++
		fmt.Printf("🔄 Aggregation cycle #%d\n", count)
		
		if err := scrapeFeeds(s); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Println("⏳ Will try again on next tick...")
		}
		
		fmt.Printf("⏳ Waiting %s until next fetch...\n\n", timeBetweenRequests)
		// Wait for next tick
		<-ticker.C
	}
}
