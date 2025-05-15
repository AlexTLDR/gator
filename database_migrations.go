package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

func ensureDatabaseMigrations(db *sql.DB) error {
	if err := ensureUsersTable(db); err != nil {
		return fmt.Errorf("ensure users table: %w", err)
	}

	if err := ensureFeedsTable(db); err != nil {
		return fmt.Errorf("ensure feeds table: %w", err)
	}

	return nil
}

func ensureUsersTable(db *sql.DB) error {
	var exists bool
	err := db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_name = 'users'
        )`,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check users table exists: %w", err)
	}

	// If table doesn't exist, create it
	if !exists {
		log.Println("Creating users table...")
		_, err = db.ExecContext(
			context.Background(),
			`CREATE TABLE users (
                id UUID PRIMARY KEY,
                created_at TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NOT NULL,
                name TEXT NOT NULL UNIQUE
            )`,
		)
		if err != nil {
			return fmt.Errorf("create users table: %w", err)
		}
	}

	return nil
}

func ensureFeedsTable(db *sql.DB) error {
	var exists bool
	err := db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_name = 'feeds'
        )`,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check feeds table exists: %w", err)
	}

	// If table doesn't exist, create it
	if !exists {
		log.Println("Creating feeds table...")
		_, err = db.ExecContext(
			context.Background(),
			`CREATE TABLE feeds (
                id UUID PRIMARY KEY,
                created_at TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NOT NULL,
                name TEXT NOT NULL,
                url TEXT NOT NULL UNIQUE,
                user_id UUID NOT NULL,
                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
            )`,
		)
		if err != nil {
			return fmt.Errorf("create feeds table: %w", err)
		}
	}

	return nil
}
