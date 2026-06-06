# Prerequisites
- PostgreSQL
- Go

# Installation

## PostgreSQL (Arch linux)

Install PostgreSQL:
```bash
sudo pacman -S postgresql
```

Initialize database cluster:
```bash
sudo -iu postgres initdb -D /var/lib/postgres/data
```

Update postgres password:
```bash
sudo passwd postgres
```

Start the service:
```bash
sudo systemctl start postgresql
```

Enter the psql shell:
```bash
sudo -u postgres psql
```

Create the database:
```sql
CREATE DATABASE gator;
```

Connect to the database:
```
\c gator
```

Set the user password:
```sql
ALTER USER postgres PASSWORD 'password';
```

Run the database migrations:
```
goose -dir sql/schema/ postgres "postgres://postgres:password@localhost:5432/gator" up
```

## Go (Arch linux)

Install Go:
```bash
sudo pacman -S go
```

## gator

Install gator:
```bash
go install https://github.com/314159otr/gator
```

## Config file

Create a `.gatorconfig.json` file in your home directory with the structure:
```json
{
    "DbURL":"postgres://postgres:password@localhost:5432/gator?sslmode=disable"
}
```

# Usage

Create a new user:
```bash
gator register user
```

Add a feed:
```bash
gator addfeed boot.dev https://www.boot.dev/blog/index.xml
```

Start the aggregator:
```bash
gator agg 30s
```

View posts:
```bash
gator browse 10
```

All commands:
- `login <username>`
- `register <username>`
- `reset`
- `users`
- `agg <time>`
- `addfeed <name> <url>`
- `feeds`
- `follow <url>`
- `following`
- `unfollow <url>`
- `browse [limit]`
