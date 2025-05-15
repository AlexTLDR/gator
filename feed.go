package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
	"time"
)

// RSSFeed represents a parsed RSS feed
type RSSFeed struct {
	XMLName     xml.Name   `xml:"rss"`
	Version     string     `xml:"version,attr"`
	Channel     RSSChannel `xml:"channel"`
}

// RSSChannel represents the channel element of an RSS feed
type RSSChannel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	PubDate     string    `xml:"pubDate"`
	LastBuildDate string  `xml:"lastBuildDate"`
	Items       []RSSItem `xml:"item"`
}

// RSSItem represents a single item in an RSS feed
type RSSItem struct {
	Title       string   `xml:"title"`
	Link        string   `xml:"link"`
	Description string   `xml:"description"`
	PubDate     string   `xml:"pubDate"`
	GUID        string   `xml:"guid"`
	Categories  []string `xml:"category"`
}

// fetchFeed retrieves and parses an RSS feed from the given URL
func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// Create a new HTTP request with the provided context
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	// Set the User-Agent header to identify our application
	req.Header.Set("User-Agent", "gator")

	// Send the request
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching feed: %w", err)
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Parse the XML
	var feed RSSFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("error parsing feed: %w", err)
	}

	// Unescape HTML entities in titles and descriptions
	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	
	for i := range feed.Channel.Items {
		feed.Channel.Items[i].Title = html.UnescapeString(feed.Channel.Items[i].Title)
		feed.Channel.Items[i].Description = html.UnescapeString(feed.Channel.Items[i].Description)
	}

	return &feed, nil
}

// handlerAgg handles the "agg" command that fetches and displays an RSS feed
func handlerAgg(s *state, cmd command) error {
	ctx := context.Background()
	
	// For now, we're just fetching a single feed
	feedURL := "https://www.wagslane.dev/index.xml"
	
	feed, err := fetchFeed(ctx, feedURL)
	if err != nil {
		return fmt.Errorf("error fetching feed: %w", err)
	}
	
	// Print the feed details
	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Link: %s\n", feed.Channel.Link)
	fmt.Printf("Items: %d\n\n", len(feed.Channel.Items))
	
	for i, item := range feed.Channel.Items {
		fmt.Printf("Item %d:\n", i+1)
		fmt.Printf("  Title: %s\n", item.Title)
		fmt.Printf("  Link: %s\n", item.Link)
		fmt.Printf("  Published: %s\n", item.PubDate)
		fmt.Printf("  Description: %s\n\n", item.Description)
	}
	
	return nil
}