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

Or you can just save a lot of time by just using Docker! create a .env file
with keys that you can find in [.env-example](.env-example) and fill to your
liking. After that you can do `docker-compose up -d` to run the containers.

If you want to centralize the logging of the services with Grafana Loki.
You need to install the docker plugin first by doing:

```bash
docker plugin install grafana/loki-docker-driver:latest --alias loki --grant-all-permissions

# restart docker service
systemctl restart docker

# verify if plugin is installed
docker plugin ls
```

You should see the installed plugin like this:

```bash
ddd2367c8693   loki:latest   Loki Logging Driver   true
```

After that you can `cd` to the `loki` directory and run:

```bash
docker-compose up -d --build
```
