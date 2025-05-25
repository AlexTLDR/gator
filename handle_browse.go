package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/AlexTLDR/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	// Set default limit
	limit := 2

	// Check if a limit was provided
	if len(cmd.Args) > 0 {
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		if limit < 1 {
			return fmt.Errorf("limit must be a positive number")
		}
	}

	// Get posts for the user
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return fmt.Errorf("error getting posts: %w", err)
	}

	// Check if any posts were found
	if len(posts) == 0 {
		fmt.Println("No posts found. Try following some feeds first!")
		return nil
	}

	// Display the posts
	fmt.Printf("Recent posts from feeds you follow (showing %d):\n\n", len(posts))
	for i, post := range posts {
		fmt.Printf("=== %d ===\n", i+1)
		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Feed: %s\n", post.FeedName)
		
		if post.PublishedAt.Valid {
			fmt.Printf("Published: %v\n", post.PublishedAt.Time.Format("January 2, 2006 15:04:05"))
		}
		
		fmt.Printf("URL: %s\n", post.Url)
		
		if post.Description.Valid && post.Description.String != "" {
			// Display a truncated description if it's too long
			desc := post.Description.String
			if len(desc) > 100 {
				desc = desc[:100] + "..."
			}
			fmt.Printf("Description: %s\n", desc)
		}
		
		fmt.Println()
	}

	// Show information about browsing more posts
	fmt.Printf("To view more posts, use: browse <limit>\n")
	
	return nil
}