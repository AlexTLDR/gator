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
