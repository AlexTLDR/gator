# Gator

An RSS feed aggreGATOR in Go! A CLI tool that allows users to:

* Add RSS feeds from across the internet to be collected
* Store the collected posts in a PostgreSQL database
* Follow and unfollow RSS feeds that other users have added
* View summaries of the aggregated posts in the terminal, with a link to the full post

## Prerequisites

To use Gator, you need to have the following installed:

* **Go** (version 1.24 or later): [Download and install from golang.org](https://golang.org/dl/)
* **PostgreSQL** (version 17 or later): [Download and install from postgresql.org](https://www.postgresql.org/download/)

## Installation

### Using go install

```bash
# Clone the repository
git clone https://github.com/AlexTLDR/gator.git

# Navigate to the project directory
cd gator

# Install the CLI
go install
```

After installation, make sure your Go bin directory is in your system PATH. You should now be able to run the `gator` command from anywhere.

## Configuration

Before using Gator, you need to set up a configuration file. Create a file named `.gatorconfig.json` in your home directory with the following content:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user_name": ""
}
```

Replace `username` and `password` with your PostgreSQL credentials. You'll need to create a database named `gator` in PostgreSQL:

```sql
CREATE DATABASE gator;
```

## Usage

### User Management

```bash
# Register a new user
gator register <username>

# Login as an existing user
gator login <username>

# View all registered users
gator users

# Reset (delete all users and their data)
gator reset
```

### Feed Management

```bash
# Add a new feed (automatically follows it)
gator addfeed "<feed_name>" "<feed_url>"

# List all feeds in the system
gator feeds

# Follow a feed by URL
gator follow <feed_url>

# Unfollow a feed
gator unfollow <feed_url>

# Show feeds you're following
gator following
```

### Content Aggregation and Browsing

```bash
# Start the aggregator (runs continuously)
# Collects posts every 30 seconds
gator agg 30s

# Browse posts from feeds you follow
gator browse       # Show default number of posts
gator browse 10    # Show up to 10 posts
```

## Tips for Using Gator

1. The aggregator (`gator agg`) runs as a continuous process. You can leave it running in one terminal while using other commands in another terminal.

2. Press Ctrl+C to stop the aggregator.

3. Try these example RSS feeds:
   - TechCrunch: https://techcrunch.com/feed/
   - Hacker News: https://news.ycombinator.com/rss
   - Boot.dev Blog: https://blog.boot.dev/index.xml

4. The `browse` command shows the most recent posts from feeds you follow. Use a higher limit to see more posts.

5. If you encounter any issues, check your PostgreSQL connection and make sure the database exists.
