package main

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexTLDR/gator/internal/database"

	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %v <name>", cmd.Name)
	}

	name := cmd.Args[0]

	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	})
	if err != nil {
		return fmt.Errorf("couldn't create user: %w", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User created successfully:")
	printUser(user)
	return nil
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: %s <name>", cmd.Name)
	}
	name := cmd.Args[0]

	_, err := s.db.GetUser(context.Background(), name)
	if err != nil {
		return fmt.Errorf("couldn't find user: %w", err)
	}

	err = s.cfg.SetUser(name)
	if err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}

	fmt.Println("User switched successfully!")
	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

func handlerUsers(s *state, cmd command) error {
	// No arguments needed for this command
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v (takes no arguments)", cmd.Name)
	}

	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't retrieve users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found in the database.")
		return nil
	}

	currentUser, err := s.cfg.GetUser()
	if err != nil {
		// If error getting current user, just print all users without marking anyone as current
		for _, user := range users {
			fmt.Printf("* %s\n", user.Name)
		}
		return nil
	}

	for _, user := range users {
		if user.Name == currentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}

func handlerReset(s *state, cmd command) error {
	// No arguments needed for this command
	if len(cmd.Args) != 0 {
		return fmt.Errorf("usage: %v (takes no arguments)", cmd.Name)
	}

	err := s.db.DeleteAllUsers(context.Background())
	if err != nil {
		return fmt.Errorf("couldn't delete users: %w", err)
	}

	fmt.Println("All users and their associated feeds have been successfully deleted!")
	return nil
}
