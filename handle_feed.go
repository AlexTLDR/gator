package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexTLDR/gator/internal/database"

	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command) error {
	// Check for correct number of arguments
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage: %v <name> <url>", cmd.Name)
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	currentUserName, err := s.cfg.GetUser()
	if err != nil {
		return fmt.Errorf("couldn't get current user: %w", err)
	}

	user, err := s.db.GetUser(context.Background(), currentUserName)
	if err != nil {
		return fmt.Errorf("couldn't find user '%s': %w", currentUserName, err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed: %w", err)
	}

	fmt.Println("Feed created successfully:")
	printFeed(feed)
	return nil
}

func printFeed(feed database.Feed) {
	fmt.Printf(" * ID:        %v\n", feed.ID)
	fmt.Printf(" * Name:      %v\n", feed.Name)
	fmt.Printf(" * URL:       %v\n", feed.Url)
	fmt.Printf(" * User ID:   %v\n", feed.UserID)
	fmt.Printf(" * Created:   %v\n", feed.CreatedAt)
	fmt.Printf(" * Updated:   %v\n", feed.UpdatedAt)
}

func handlerFeeds(s *state, cmd command) error {
	// No arguments needed for this command
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v (takes no arguments)", cmd.Name)
	}

	// Get all feeds with user information
	feeds, err := s.db.GetFeedsWithUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve feeds: %w", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found in the database.")
		return nil
	}

	fmt.Println("Feeds:")
	for i, feed := range feeds {
		fmt.Printf("%d. %s\n", i+1, feed.Name)
		fmt.Printf("   URL:  %s\n", feed.Url)
		fmt.Printf("   User: %s\n\n", feed.UserName)
	}

	return nil
}
