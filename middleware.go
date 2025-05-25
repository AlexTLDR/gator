package main

import (
	"context"
	"fmt"

	"github.com/AlexTLDR/gator/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		// Get current user name from config
		currentUserName, err := s.cfg.GetUser()
		if err != nil {
			return fmt.Errorf("couldn't get current user: %w", err)
		}

		if currentUserName == "" {
			return fmt.Errorf("you must be logged in to use this command. Use 'login <username>' or 'register <username>' first")
		}

		// Find the user in the database
		user, err := s.db.GetUser(context.Background(), currentUserName)
		if err != nil {
			return fmt.Errorf("couldn't find user '%s': %w", currentUserName, err)
		}

		// Call the original handler with the user
		return handler(s, cmd, user)
	}
}
