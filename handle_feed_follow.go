package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/AlexTLDR/gator/internal/database"

	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	// Check for correct number of arguments
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <url>", cmd.Name)
	}

	url := cmd.Args[0]

	// Find feed by URL
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("couldn't find feed with URL '%s': %w", url, err)
	}

	// Check if user is already following this feed
	_, err = s.db.GetFeedFollowByUserAndFeed(context.Background(), database.GetFeedFollowByUserAndFeedParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	
	// If no error, the feed follow already exists
	if err == nil {
		fmt.Printf("You are already following '%s'\n", feed.Name)
		return nil
	}
	
	// If the error is not "no rows", then it's an unexpected error
	if err != sql.ErrNoRows {
		return fmt.Errorf("error checking existing feed follow: %w", err)
	}

	// Create feed follow
	now := time.Now().UTC()
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("couldn't create feed follow: %w", err)
	}

	fmt.Printf("You (%s) are now following '%s'\n", feedFollow.UserName, feedFollow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	// Check for correct number of arguments
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v (takes no arguments)", cmd.Name)
	}

	// Get feed follows for user
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("couldn't get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Printf("User '%s' is not following any feeds.\n", user.Name)
		return nil
	}

	fmt.Printf("Feeds followed by '%s':\n", user.Name)
	fmt.Println("-----------------------------")
	for i, ff := range feedFollows {
		fmt.Printf("%d. %s\n", i+1, ff.FeedName)
		fmt.Printf("   Started following: %s\n", ff.CreatedAt.Format("Jan 02, 2006 15:04:05"))
	}
	fmt.Printf("\nTotal feeds followed: %d\n", len(feedFollows))

	return nil
}

func printFeedFollow(feedFollow database.CreateFeedFollowRow) {
	fmt.Printf(" * ID:        %v\n", feedFollow.ID)
	fmt.Printf(" * User:      %v\n", feedFollow.UserName)
	fmt.Printf(" * Feed:      %v\n", feedFollow.FeedName)
	fmt.Printf(" * Created:   %v\n", feedFollow.CreatedAt)
	fmt.Printf(" * Updated:   %v\n", feedFollow.UpdatedAt)
}