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

	if err := ensureFeedFollowsTable(db); err != nil {
		return fmt.Errorf("ensure feed_follows table: %w", err)
	}

	if err := ensureFeedsLastFetchedColumn(db); err != nil {
		return fmt.Errorf("ensure feeds last_fetched_at column: %w", err)
	}

	if err := ensurePostsTable(db); err != nil {
		return fmt.Errorf("ensure posts table: %w", err)
	}

	return nil
}

func ensurePostsTable(db *sql.DB) error {
	var exists bool
	err := db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_name = 'posts'
        )`,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check posts table exists: %w", err)
	}

	// If table doesn't exist, create it
	if !exists {
		log.Println("Creating posts table...")
		_, err = db.ExecContext(
			context.Background(),
			`CREATE TABLE posts (
                id UUID PRIMARY KEY,
                created_at TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NOT NULL,
                title TEXT NOT NULL,
                url TEXT NOT NULL UNIQUE,
                description TEXT,
                published_at TIMESTAMP,
                feed_id UUID NOT NULL,
                FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
            )`,
		)
		if err != nil {
			return fmt.Errorf("create posts table: %w", err)
		}
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

func ensureFeedsLastFetchedColumn(db *sql.DB) error {
	var exists bool
	err := db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS (
			SELECT FROM information_schema.columns
			WHERE table_name = 'feeds' AND column_name = 'last_fetched_at'
		)`,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check feeds last_fetched_at column exists: %w", err)
	}

	// If column doesn't exist, create it
	if !exists {
		log.Println("Adding last_fetched_at column to feeds table...")
		_, err = db.ExecContext(
			context.Background(),
			`ALTER TABLE feeds
			ADD COLUMN last_fetched_at TIMESTAMP NULL`,
		)
		if err != nil {
			return fmt.Errorf("add last_fetched_at column to feeds table: %w", err)
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

func ensureFeedFollowsTable(db *sql.DB) error {
	var exists bool
	err := db.QueryRowContext(
		context.Background(),
		`SELECT EXISTS (
            SELECT FROM information_schema.tables
            WHERE table_name = 'feed_follows'
        )`,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("check feed_follows table exists: %w", err)
	}

	// If table doesn't exist, create it
	if !exists {
		log.Println("Creating feed_follows table...")
		_, err = db.ExecContext(
			context.Background(),
			`CREATE TABLE feed_follows (
                id UUID PRIMARY KEY,
                created_at TIMESTAMP NOT NULL,
                updated_at TIMESTAMP NOT NULL,
                user_id UUID NOT NULL,
                feed_id UUID NOT NULL,
                UNIQUE(user_id, feed_id),
                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
                FOREIGN KEY (feed_id) REFERENCES feeds(id) ON DELETE CASCADE
            )`,
		)
		if err != nil {
			return fmt.Errorf("create feed_follows table: %w", err)
		}
	}

	return nil
}
