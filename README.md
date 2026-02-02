# dwarf

Corporate production bullshit stuff that I need to reinforce to my learnings.

## Running Locally

### Prerequisites

Make sure you already installed `golang-migration` and `sqlc` by:

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && \
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Ensure the PostgreSQL is up and running locally and the target database
already exists.

### Setup

Run database migrations:

```bash
make migrate-up DATABASE_URI=...
```

Generate database schema and sqlc code for the server:

```bash
make schema
```

### Run the server

```bash
make run-server \ 
    DATABASE_USER=... \
    DATABASE_PASSWORD=... \
    DATABASE_PORT=... \
    DATABASE_NAME=...
```
