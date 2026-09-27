# Gator

Gator is a CLI RSS feed aggregator written in Go. It allows users to register accounts, follow RSS feeds, aggregate posts, and browse content from feeds they follow.

## Requirements

Before running Gator, make sure you have the following installed:

- Go (1.23+ recommended)
- PostgreSQL

Verify installations:

```bash
go version
psql --version
```

## Installation

Install the CLI using:

```bash
go install github.com/mohammad-rizwan-hussain/gator@latest
```

Or install from a local checkout:

```bash
go install .
```

Ensure that Go's bin directory is on your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Verify installation:

```bash
gator
```

## Database Setup

Create a PostgreSQL database:

```sql
CREATE DATABASE gator;
```

Run the database migrations:

```bash
goose postgres "<connection_string>" up
```

Example connection string:

```text
postgres://username:password@localhost:5432/gator?sslmode=disable
```

## Configuration

Create a configuration file at:

```text
~/.gatorconfig.json
```

Example:

```json
{
  "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
  "current_user": ""
}
```

## Usage

### Register a User

```bash
gator register kahya
```

### Login

```bash
gator login kahya
```

### View Users

```bash
gator users
```

### Add a Feed

```bash
gator addfeed "TechCrunch" https://techcrunch.com/feed/
```

### List Feeds

```bash
gator feeds
```

### Follow a Feed

```bash
gator follow https://techcrunch.com/feed/
```

### View Followed Feeds

```bash
gator following
```

### Start Feed Aggregation

Fetch feeds every 30 seconds:

```bash
gator agg 30s
```

### Browse Posts

Show the 2 most recent posts:

```bash
gator browse
```

Show 10 posts:

```bash
gator browse 10
```

## Available Commands

| Command | Description |
|----------|-------------|
| register | Register a new user |
| login | Log in as a user |
| users | List all users |
| addfeed | Add a new RSS feed |
| feeds | List all feeds |
| follow | Follow a feed |
| following | List followed feeds |
| agg | Continuously fetch feed content |
| browse | Browse posts from followed feeds |
| reset | Reset the database |

## Project Structure

```text
.
├── internal/
│   ├── config/
│   ├── database/
│   ├── handlers/
│   ├── rss/
│   └── cmd/
├── sql/
├── migrations/
└── main.go
```

## License

This project was created as part of the Boot.dev curriculum.
